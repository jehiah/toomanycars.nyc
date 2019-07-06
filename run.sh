#!/bin/bash

# https://dev.socrata.com/foundry/data.cityofnewyork.us/p2mh-mrfv
# API Docs ^^

FIELDS=license_nbr,license_status,license_creation_date,lic_expir_dd,industry,business_name,business_name_2,address_building,address_street_name,address_borough,detail_2
# >>> urllib.quote("industry='Garage and Parking Lot' OR industry='Parking Lot' OR industry=Garage")
WHERE="industry%3D%27Garage%20and%20Parking%20Lot%27%20OR%20industry%3D%27Parking%20Lot%27%20OR%20industry%3DGarage"
WHERE="industry='Garage%20and%20Parking%20Lot'%20OR%20industry='Parking%20Lot'%20OR%20industry='Garage'"
echo "https://data.cityofnewyork.us/resource/p2mh-mrfv.json?\$where=${WHERE}&\$select=${FIELDS}&\$order=license_creation_date%20ASC\$limit=5000"

curl "https://data.cityofnewyork.us/resource/p2mh-mrfv.json?\$where=${WHERE}&\$select=${FIELDS}&\$order=license_creation_date%20ASC&\$limit=5000" --silent  > data/dca_licenses.json
