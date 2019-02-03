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
parser.add_argument("-s", "--stage", help="Värde för STAGE-etiketten för den peer connection du vill ta bort.")
args = parser.parse_args()

if args.region == None:
    print(parser.print_help())
    exit(1)

print("Hämtar aktiva peer connections")

pcs = json.loads(run_command(["aws", "ec2", "describe-vpc-peering-connections", "--region", args.region, "--filters", "Name=status-code,Values=active"]))['VpcPeeringConnections']

for index, pc in enumerate(pcs):
    print("-" * 30)
    print("Nr %d: %s" % (index, pc['VpcPeeringConnectionId']))
    print("-" * 30)

pc_index = int(input("Ange vilken peer connection du vill ta bort: "))
pc = pcs[pc_index]

print("Tar bort routes från route tables som använder den peer connection som ska tas bort...")

route_tables = json.loads(run_command(["aws", "ec2", "describe-route-tables", "--region", args.region, "--filters", "Name=route.vpc-peering-connection-id,Values=" + pc['VpcPeeringConnectionId']]))['RouteTables']

print("Hittade %d route tables." % (len(route_tables)))

for route_table in route_tables:
    for route in [r for r in route_table['Routes'] if 'VpcPeeringConnectionId' in r and r['VpcPeeringConnectionId'] == pc['VpcPeeringConnectionId']]:
        print("Tar bort route med cidr block %s i route table %s" % (route_table['RouteTableId'], route['DestinationCidrBlock']))
        run_command(["aws", "ec2", "delete-route", "--region", args.region, "--route-table-id", route_table['RouteTableId'], "--destination-cidr-block", route['DestinationCidrBlock']])

print("Hämtar security groups...")
security_groups = json.loads(run_command(["aws", "ec2", "describe-security-groups", "--region", args.region, "--filters", "Name=tag:aws:cloudformation:logical-id,Values=SheltersVPCSecurityGroup", "Name=tag:STAGE,Values=" + args.stage]))['SecurityGroups']

print("Hittade %d security groups" % (len(security_groups)))

for group in security_groups:
    print("-" * 30)
    print("GroupId: %s" % (group['GroupId']))
    print("GroupName: %s" % (group['GroupName']))
    for perm in group['IpPermissions']:
        print("Protokoll: %s" % (perm['IpProtocol']))
        print("Från port: %s" % (perm['FromPort']))
        print("Till port: %s" % (perm['ToPort']))
        for ip in perm['IpRanges']:
            print("IP-omfång: %s" % (ip['CidrIp']))
        ans = input("Vill du ta bort ovanstående ingress? [y/n] ")

        if ans != 'y':
            print("Hoppar över.")
            continue

        print("Tar bort...")
        for ip in perm['IpRanges']:
            run_command(["aws", "ec2", "revoke-security-group-ingress", "--region", args.region, "--group-id", group['GroupId'], "--cidr", ip['CidrIp'], "--protocol", perm['IpProtocol'], "--port", str(perm['FromPort'])])
        print("Ingress borttagen.")
    print("-" * 30)

print("Tar bort peering connection...")

run_command(["aws", "ec2", "delete-vpc-peering-connection", "--region", args.region, "--vpc-peering-connection-id", pc['VpcPeeringConnectionId']])

print("\nKlart.")
