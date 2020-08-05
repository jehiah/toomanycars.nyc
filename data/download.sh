#!/bin/bash


#################
# DCA LICENSE DATA
#################
# https://dev.socrata.com/foundry/data.cityofnewyork.us/p2mh-mrfv
# API Docs ^^
DATASET=p2mh-mrfv
DATASET=w7w3-xahh

FIELDS=license_nbr,license_status,license_creation_date,lic_expir_dd,industry,business_name,business_name_2,address_building,address_street_name,address_borough,detail_2
WHERE="industry='Garage%20and%20Parking%20Lot'%20OR%20industry='Parking%20Lot'%20OR%20industry='Garage'"

echo "downloading dca_licenses.json"
curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$where=${WHERE}&\$select=${FIELDS}&\$order=license_nbr%20ASC&\$limit=9000" --silent  > dca_licenses_tmp.json

# reformat to one line per record
echo -n "[" > dca_licenses.json
cat dca_licenses_tmp.json | jq -c '.[]' >> dca_licenses.json
gsed -i -e 's/^{/,{/g' dca_licenses.json
echo "]" >> dca_licenses.json
rm dca_licenses_tmp.json


###################
# Parking Lot Map Data
###################

# https://data.cityofnewyork.us/City-Government/Parking-Lot/h7zy-iq3d
DATASET=h7zy-iq3d
DATASET=7cgt-uhhz
FIELDS="source_id,status,shape_leng,shape_area"
if [ ! -f DOITT_planimetrics_parking_lot.json ]; then
	echo "downloading DOITT_planimetrics_parking_lot.json"
	curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$select=${FIELDS}&\$limit=200000" --silent  > DOITT_planimetrics_parking_lot.json
fi
if [ ! -f DOITT_planimetrics_parking_lot.geojson ]; then
	echo "downloading DOITT_planimetrics_parking_lot.geojson"
	curl "https://data.cityofnewyork.us/resource/${DATASET}.geojson?\$select=${FIELDS},the_geom&\$limit=200000" --silent  > DOITT_planimetrics_parking_lot.geojson
fi


#################
# Building Footprint
# Subtype: Garage
###############

# https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh
# https://dev.socrata.com/foundry/data.cityofnewyork.us/6kx9-25sv
DATASET=nqwf-w8eh
DATASET=xra2-rhxp
FIELDS="doitt_id,bin,feat_code,shape_area"
WHERE="feat_code=5110"
if [ ! -f DOITT_planimetrics_building_garages.json ]; then
	echo "downloading DOITT_planimetrics_building_garages.json"
	curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$where=${WHERE}&\$select=${FIELDS}&\$limit=500000" --silent  > DOITT_planimetrics_building_garages.json
fi
if [ ! -f DOITT_planimetrics_building_garages.geojson ]; then
	echo "downloading DOITT_planimetrics_building_garages.geojson"
	curl "https://data.cityofnewyork.us/resource/${DATASET}.geojson?\$where=${WHERE}&\$select=${FIELDS},the_geom&\$limit=500000" --silent  > DOITT_planimetrics_building_garages.geojson
fi

# curl -v 'https://data.cityofnewyork.us/api/views/metadata/v1/h7zy-iq3d'

#######
# borough geojson
# 
