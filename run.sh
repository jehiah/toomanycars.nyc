#!/bin/bash


#################
# DCA LICENSE DATA
#################
# https://dev.socrata.com/foundry/data.cityofnewyork.us/p2mh-mrfv
# API Docs ^^
DATASET=p2mh-mrfv
DATASET=w7w3-xahh

FIELDS=license_nbr,license_status,license_creation_date,lic_expir_dd,industry,business_name,business_name_2,address_building,address_street_name,address_borough,detail_2
# >>> urllib.quote("industry='Garage and Parking Lot' OR industry='Parking Lot' OR industry=Garage")
WHERE="industry%3D%27Garage%20and%20Parking%20Lot%27%20OR%20industry%3D%27Parking%20Lot%27%20OR%20industry%3DGarage"
WHERE="industry='Garage%20and%20Parking%20Lot'%20OR%20industry='Parking%20Lot'%20OR%20industry='Garage'"
echo "https://data.cityofnewyork.us/resource/${DATASET}.json?\$where=${WHERE}&\$select=${FIELDS}&\$order=license_creation_date%20ASC\$limit=5000"

curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$where=${WHERE}&\$select=${FIELDS}&\$order=license_creation_date%20ASC&\$limit=5000" --silent  > data/dca_licenses_tmp.json

# reformat to one line per record
echo -n "[" > data/dca_licenses.json
cat data/dca_licenses_tmp.json | jq -c '.[]' >> data/dca_licenses.json
gsed -i -e 's/^{/,{/g' data/dca_licenses.json
echo "]" >> data/dca_licenses.json
rm data/dca_licenses_tmp.json


###################
# Parking Lot Map Data
###################

DATASET=h7zy-iq3d
DATASET=7cgt-uhhz
FIELDS="source_id,status,shape_leng,shape_area"
curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$select=${FIELDS}&\$limit=200000" --silent  > data/DOITT_planimetrics_parking_lot.json


#################
# Building Footprint
# Subtype: Garage
###############

# https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh
# https://dev.socrata.com/foundry/data.cityofnewyork.us/6kx9-25sv
DATASET=nqwf-w8eh
DATASET=6kx9-25sv
FIELDS="doitt_id,bin,feat_code,shape_area"
WHERE="feat_code=5110"
curl "https://data.cityofnewyork.us/resource/${DATASET}.json?\$where=${WHERE}&\$select=${FIELDS}&\$limit=500000" --silent  > data/DOITT_planimetrics_building_garages.json
