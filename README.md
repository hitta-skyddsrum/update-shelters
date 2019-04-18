# update-shelters

## check-for-updates
`bin/check` downloads the latest shapefile based on the URL given in the [XML feed](https://gis-services.metria.se/msbfeed/skyddsrum.xml)

## parse-shapefile
`bin/parse-shapefile shelters.zip` will create a JSON file include all the shelters in given shapefile.
```
$ bin/parse-shapefile

Usage of bin/parse-shapefile:
  -list-fields
    	List all fields in shapefile
  -show-example
    	Show an example shape from the shapefile
```

## Deploy
```
aws cloudformation package --template-file cloudformation.yml --s3-bucket update-shelters-deploy-bucket --output-template cf-packaged.yml
aws cloudformation deploy --stack-name update-shelters --template-file cf-packaged.yml --capabilities CAPABILITY_IAM --parameter-overrides BucketName= DbMasterUser=DbMasterPassword=
```

## Read shapefile manually
1. Open QGIS
1. Select Layer > Add Layer > Add Vector Layer
1. Select Layer > Open Attribute Table
