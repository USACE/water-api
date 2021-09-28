'''
Cleanup uses Shapely's make_valid()
https://shapely.readthedocs.io/en/latest/manual.html#validation.make_valid

Returns a valid representation of the geometry, if it is invalid. If it is 
valid, the input geometry will be returned. In many cases, in order to 
create a valid geometry, the input geometry must be split into multiple 
parts or multiple geometries. If the geometry must be split into multiple 
parts of the same geometry type, then a multi-part geometry (e.g. a
MultiPolygon) will be returned. if the geometry must be split into multiple 
parts of different types, then a GeometryCollection will be returned.
'''
from shapely.validation import make_valid
from shapely.geometry import shape

'''
Takes in a Collection of features.
Returns a list of Shapely Geometry objects
'''
def process(input_shapes):

    shapes = []
    invalid_shapes = []
    for shp in input_shapes:
        if not shp.is_valid:
            # Shape is not valid, apply fix
            invalid_shapes.append(shp)
            valid = make_valid(shp.buffer(0.5))
            shapes.append(valid)            
        else:
            shapes.append(shape(shp.buffer(0.5)))

    print(f'{len(invalid_shapes)} of {len(shapes)} shapes were invalid and corrected.')   
    return shapes                                   