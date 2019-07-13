
https://data.cityofnewyork.us/dataset/Parking-Facilities1/ktgb-gg65
    > Parking Facilities1
    > A list of parking facilities that have a current DCA license as of the run date.
    > Data Last Updated August 20, 2011
    *** Has parking spaces
    

https://data.cityofnewyork.us/Transportation/Municipal-Parking-Facilities-Manhattan-/i8d5-5ciu
    > same for each borough
    > updated September 10, 2018


--

https://data.cityofnewyork.us/Transportation/Parking-Meters-ParkNYC-Blockfaces/s7zi-dgdx
Parking meter block faces
    https://data.cityofnewyork.us/Transportation/Parking-Meters-ParkNYC-Blockfaces/s7zi-dgdx
    

https://data.cityofnewyork.us/Business/License-Revocations-Suspensions-Surrenders-and-Rei/rpeq-j89e/data
    DCA License revocation, etc
        industry in [Garage, Parking Lot, Garage and Parking Lot]

DCA License Check
    https://a858-elpaca.nyc.gov/CitizenAccess/

https://data.cityofnewyork.us/Business/Legally-Operating-Businesses/w7w3-xahh/data
DCA Legally Operating Businesses (updated regularly)
    > Garage and Parking Lot
    > Garage
    > Parking Lot
        > Detail Column
    
    LicenseNumber,LicenseCreationDate,LicenseExpireDate,Industry,BusinessName,BusinessName2,AddressBuilding,AddressStreetName,AddressBorough,Detail
    

Parking Lot Maps https://data.cityofnewyork.us/City-Government/Parking-Lot/h7zy-iq3d
https://data.cityofnewyork.us/resource/7cgt-uhhz.json
    source_id, status, shape_area
    20965sqft/27=776sq per space

https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh
Building Footprint
    subtype: garage: 5110 (aka 'feat_code')

NYC BIN lookup
http://a030-goat.nyc.gov/goat/FunctionBN?bin=null

Curb Cut - New
https://data.cityofnewyork.us/Housing-Development/DOB-NOW-Build-Approved-Permits/rbx6-tga4/data
    > filter WorkType = "Curb Cut"
    > 

https://www1.nyc.gov/site/buildings/industry/project-requirements-design-professional-curb-cuts.page
T-001.00	

BC=Building Code
driveway size and location 	
2014 BC 406.7.7
https://www1.nyc.gov/assets/buildings/apps/pdf_viewer/viewer.html?file=2014CC_BC_Chapter_4_Special_Detailed_Requirements.pdf&section=conscode_2014

# NYC Buildings - Certificate of Occupancy

Private garages and carports, as defined by this section, shall be classified as Group U occu
 Certificate of Occupancy>
     > garage 
      Certificate of Occupancy shall indicate the maximum number of vehicles to be accommodated and the type of vehicle, whether private passenger or commercial, to be stored. An application for or including an open parking lot shall be accompanied by a plan exhibiting the following:
      
https://data.cityofnewyork.us/Housing-Development/DOB-Certificate-Of-Occupancy/bs8b-p36w/data
    >
    BIN 1083785  http://a810-bisweb.nyc.gov/bisweb/CofoDocumentContentServlet?passjobnumber=null&cofomatadata1=cofo&cofomatadata2=M&cofomatadata3=120&cofomatadata4=526000&cofomatadata5=120526911.PDF
    > PARKING GARAGE FOR 205 MOTOR VEHICLES
    
    BIN: 4100767 - CO 421412359T001
    > 27 OUTDOOR ACCESSORY PARKING SPACES AND 2 LOADING BERTHS
    
    BIN: 4044890 - CO 401066047T001
    > 58 ACCESSORY OPEN PARKING SPACES OG FULLY ATTENDED
    
    BIN: 4566130 - CO 402159123F
    > 9 OFF STREET OPEN PARKING SPACES
    
Buildings built before 1938 aren’t required to have a Certificate of Occupancy – unless later alterations changed its use, egress or occupancy. If you require proof of a building’s legal use – and it’s exempt from the CO requirement – contact the Department’s borough office where the property is located to request a Letter of No Objection.


Certificate of Occupancy Lookup: http://a810-bisweb.nyc.gov/bisweb/bispi00.jsp
    > BIN lookup - http://a810-bisweb.nyc.gov/bisweb/PropertyProfileOverviewServlet?bin=4096312&go4=+GO+&requestid=0
    > Certificate of Occupancy - http://a810-bisweb.nyc.gov/bisweb/COsByLocationServlet?requestid=1&allbin=4096312
        

Planimetrics base image
https://github.com/CityOfNewYork/nyc-geo-metadata
https://github.com/CityOfNewYork/nyc-geo-metadata/blob/master/Metadata/Metadata_AerialImagery.md


https://medium.com/the-downlinq/car-localization-and-counting-with-overhead-imagery-an-interactive-exploration-9d5a029a596b
https://medium.com/the-downlinq/yolt-arxiv-paper-and-code-release-8b30d40d095b
> https://github.com/CosmiQ/simrdwn


# Census Data

https://factfinder.census.gov/faces/tableservices/jsf/pages/productview.xhtml?pid=ACS_15_5YR_B08201&prodType=table

1981663 cars in the 5 boroughs "available for use" to households. 2013-2017 5yr estimate

# NY DMV Data

https://dmv.ny.gov/about-dmv/statistical-summaries

1,923,041 standard vehicles, 75,593 commercial, 100k taxi. 2017 data

# PA NYNJ Traffic Volumes

https://corpinfo.panynj.gov/pages/pa-traffic-volume/

50k Lincolin [EB], 43k Holland, 142k GWB

# MTA Bridge Vehicle Counts

http://web.mta.info/bandt/html/btintro.html

868k each weekday

# DOT Bridges Vehicle Counts

Ed Koch 170k, Brooklyn 124k, Williamsburg 111k, Manhattan 75k
