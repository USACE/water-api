# load csv module
import csv

param_map = {
    "00010": "Temp-Water",
    "00011": "Temp-Water",
    "00045": "Precip",
    "00061": "Flow",
}

# open file for reading
with open("lrn_usgs_locations.csv") as csvDataFile:

    # read file as csv file
    csvReader = csv.reader(csvDataFile)

    # for every row, print the row
    for row in csvReader:
        print(row[0], row[5])
