# http://pcjericks.github.io/py-gdalogr-cookbook/
import sys, os
try:
    from osgeo import ogr, osr, gdal
except:
    sys.exit('ERROR: cannot find GDAL/OGR modules')

# Enable GDAL/OGR exceptions
gdal.UseExceptions()


# daShapefile = "/test_shapes/SubbasinsAEA.shp"

# driver = ogr.GetDriverByName('ESRI Shapefile')

# dataSource = driver.Open(daShapefile, 0) # 0 means read-only. 1 means writeable.

# # Check to see if shapefile is found.
# if dataSource is None:
#     print ('Could not open %s' % (daShapefile))
# else:
#     print('Opened %s' % (daShapefile))
#     layer = dataSource.GetLayer()
#     featureCount = layer.GetFeatureCount()
#     print("Number of features in %s: %d" % (os.path.basename(daShapefile),featureCount))

file = ogr.Open("/test_shapes/SubbasinsAEA.shp")
shape = file.GetLayer(0)

layer = file.GetLayer()
print(f"Feature Count: {len(layer)}")

# for feature in layer:
#     # print feature.GetField("STATE_NAME")
#     print(feature.GetLayerDefn().GetName())
# layer.ResetReading()

# #first feature of the shapefile
# feature = shape.GetFeature(0)
# first = feature.ExportToJson()
# # print(first) # (GeoJSON format)