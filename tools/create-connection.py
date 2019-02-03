import argparse
import json
import subprocess
import sys

def run_command(cmd):
    cmd_output = subprocess.run(cmd, capture_output=True, text=True)
    if cmd_output.stderr:
        print(cmd_output.stderr)
        exit(1)
    return cmd_output.stdout

parser = argparse.ArgumentParser()
parser.add_argument("-r", "--region", help="AWS region")
parser.add_argument("-s", "--stage", help="Värde för Stage-etikett för den databas du vill ansluta till.")
args = parser.parse_args()

if args.region == None or args.stage == None:
    print(parser.print_help())
    exit(1)

print("Hämtar initial data...")

bastion_vpc_stack = json.loads(run_command(["aws", "cloudformation", "describe-stack-resource", "--region", args.region, "--stack-name", "shelters-bastion", "--logical-resource-id", "VPCStack"]))['StackResourceDetail']
bastion_vpc_stack_name = bastion_vpc_stack['PhysicalResourceId'].split(":stack/")[1].split("/")[0]
bastion_resources = json.loads(run_command(["aws", "cloudformation", "list-stack-resources", "--region", args.region, "--stack-name", bastion_vpc_stack_name]))['StackResourceSummaries']
bastion_vpc_id = [r for r in bastion_resources if r['LogicalResourceId'] == 'VPC'][0]['PhysicalResourceId']

print("Hämtar tillgängliga RDS-databaser...")

cmd = run_command(["aws", "rds", "describe-db-instances", "--region", args.region])

dbs = json.loads(cmd)['DBInstances']

for index, db in enumerate(dbs):
    cmd = run_command(["aws", "rds", "list-tags-for-resource", "--resource-name", db['DBInstanceArn'], "--region", args.region])

    tags = json.loads(cmd)['TagList']
    stage_tags = [t for t in tags if t['Key'] == 'STAGE']
    if len(stage_tags) == 0 or stage_tags[0]['Value'] != args.stage:
        continue

    print("-" * 30)
    print("Nr %d. %s" % (index, db['DBInstanceIdentifier']))
    print ('%-45s %s' % ("Maskin", db['Engine']))

    for tag in tags:
        print ('%-45s %s' % (tag['Key'], tag['Value']))
    print("-" * 30)
    print("\n")

db_index = int(input("Ange vilken databas du vill ansluta till: "))
db_address = dbs[db_index]['Endpoint']['Address']
db_port = dbs[db_index]['Endpoint']['Port']
db_vpc_id = dbs[db_index]['DBSubnetGroup']['VpcId']
db_security_group_id = [d for d in dbs[db_index]['VpcSecurityGroups'] if d['Status'] == 'active'][0]['VpcSecurityGroupId']

accepter_vpc = db_vpc_id
requester_vpc = bastion_vpc_id

print("Skapar peer connection från %s till %s..." % (requester_vpc, accepter_vpc))

cmd = run_command(["aws", "ec2", "create-vpc-peering-connection", "--peer-vpc-id", accepter_vpc , "--vpc-id", requester_vpc, "--region", args.region])

print("Peer connection skapad.")

while 1:
    cmd = run_command(["aws", "ec2", "describe-vpc-peering-connections", "--region", args.region])

    pcs = [pc for pc in json.loads(cmd)['VpcPeeringConnections'] if pc['Status']['Code'] == 'pending-acceptance' and pc['AccepterVpcInfo']['VpcId'] == accepter_vpc and pc['RequesterVpcInfo']['VpcId'] == requester_vpc]

    if len(pcs) > 0:
        break

    print("Väntar...")

if len(pcs) > 1:
    print("Hittade ett oväntat antal avvaktande peer connections: " + len(pcs) + ". Avslutar.")
    exit(1)

print("Accepterar peer connection...")

pc = pcs[0]
pc_id = pc['VpcPeeringConnectionId']

cmd = run_command(["aws", "ec2", "accept-vpc-peering-connection", "--vpc-peering-connection-id", pc_id, "--region", args.region])

json_pc = json.loads(cmd)['VpcPeeringConnection']
requester_cidr_block = json_pc['RequesterVpcInfo']['CidrBlock']
accepter_cidr_block = json_pc['AccepterVpcInfo']['CidrBlock']

cmd = run_command(["aws", "ec2", "describe-route-tables","--region", args.region, "--filters", "Name=tag:aws:cloudformation:logical-id,Values=PublicSubnetRouteTable,SheltersRouteTable", "Name=vpc-id,Values=" + bastion_vpc_id + "," + db_vpc_id]);

route_tables = json.loads(cmd)['RouteTables']

if len(route_tables) != 2:
    print("Hittade ett oväntat antal route tables: " + len(route_tables) + ". Avslutar.")
    exit(1)

for route_table in route_tables:
    if route_table['VpcId'] == accepter_vpc:
        cidr_block = requester_cidr_block
    elif route_table['VpcId'] == requester_vpc:
        cidr_block = accepter_cidr_block
    else:
        print("Kunde inte hitta skapa routes för VPC " + route_table['VpcId'] + ".")
        exti(1)

    print("Skapar route för VPC " + route_table['VpcId'] + "...")
    cmd = run_command(["aws", "ec2", "create-route", "--region", args.region, "--route-table-id", route_table['RouteTableId'], "--vpc-peering-connection-id", pc_id, "--destination-cidr-block", cidr_block])

    print("Skapade route för VPC " + route_table['VpcId'])

print("Öppnar port %d i DB-instansens security group..." % (db_port))

run_command(["aws", "ec2", "authorize-security-group-ingress", "--group-id", db_security_group_id, "--protocol", "tcp", "--port", str(db_port), "--cidr", requester_cidr_block, "--region", args.region])

print("-" * 30)
print("Skapande av SSH-tunnel lyckades.")
print("-" * 30)

# Write a summary to help the user
print("Hämtar metadata för att skriva ut sammanfattning...")

bastion_stack_id = json.loads(run_command(["aws", "cloudformation", "describe-stack-resource", "--region", args.region, "--stack-name", "shelters-bastion", "--logical-resource-id", "BastionStack"]))['StackResourceDetail']['PhysicalResourceId'].split(":stack/")[1].split("/")[0]
auto_scaling_group_name = [r for r in json.loads(run_command(["aws", "cloudformation", "list-stack-resources", "--region", args.region, "--stack-name", bastion_stack_id]))['StackResourceSummaries'] if r['ResourceType'] == 'AWS::AutoScaling::AutoScalingGroup'][0]['PhysicalResourceId']

instance_summaries = json.loads(run_command(["aws", "autoscaling", "describe-auto-scaling-groups", "--region", args.region, "--auto-scaling-group-name", auto_scaling_group_name]))['AutoScalingGroups'][0]['Instances']
print("Hittade %d instans(er)." % (len(instance_summaries)))

instances = json.loads(run_command(["aws", "ec2", "describe-instances", "--region", args.region, "--instance-ids", ",".join([i['InstanceId'] for i in instance_summaries])]))['Reservations'][0]['Instances'] 

for instance in instances:
    public_ip = instance['PublicIpAddress']
    key_name = instance['KeyName']
    print("Använd IP %s med certifikat %s som SSH-tunnel för att ansluta till %s." % (public_ip, key_name, db_address))

print("\nKlart!")
exit(0)
