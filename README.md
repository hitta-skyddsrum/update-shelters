# update-shelters

## check-for-updates
`bin/check` downloads the latest shapefile based on the URL given in the [XML feed](https://gis-services.metria.se/msbfeed/skyddsrum.xml)

## parse-shapefile
```
$ bin/parse-shapefile

Usage of bin/parse-shapefile:
  -list-fields
    	List all fields in shapefile
  -show-example
    	Show an example shape from the shapefile
```
### Output GeoJSON
`bin/parse-shapefile -geojson shelters.zip` will create a JSON file include all the shelters in given shapefile.
### Output CSV
`bin/parse-shapefile -csv shelters.zip` will create a CSV file include all the shelters in given shapefile.

## Shapefile to SQL dump
```
bin/check
unzip skyddsrum.zip
SHELTERS_SHAPEFILE=$(ls | grep "*.shp$") docker-compose up -d gdal
docker-compose gdal exec ogr2ogr -f MySQL "MySQL:hitta_skyddsrum,host=mysql,user=root,password=hitta_skyddsrum" /usr/src/shelters.shp -nln hitta_skyddsrum -overwrite -lco engine=MYISAM
```

## Read shapefile manually
1. Open QGIS
1. Select Layer > Add Layer > Add Vector Layer
1. Select Layer > Open Attribute Table
