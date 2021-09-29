'''
Cleanup uses Shapely and pyproj
https://shapely.readthedocs.io/en/stable/manual.html?highlight=transform#other-transformations
Applies func to all coordinates of geom and returns a new geometry of the same type from the transformed 
coordinates.  func maps x, y, and optionally z to output xp, yp, zp. The input parameters may be 
iterable types like lists or arrays or single values. The output shall be of the same type: scalars 
in, scalars out; lists in, lists out.

transform tries to determine which kind of function was passed in by calling func first with n 
iterables of coordinates, where n is the dimensionality of the input geometry. If func raises a 
TypeError when called with iterables as arguments, then it will instead call func on each individual 
coordinate in the geometry.
'''
from shapely.ops import transform 
from pyproj import Proj, Transformer, CRS

def process(shape, src_crs, dst_epsg=4326):

    # Credit: https://gis.stackexchange.com/questions/127427/transforming-shapely-polygon-and-multipolygon-objects
    project = Transformer.from_proj(
        Proj(src_crs), # source coordinate system
        Proj(CRS(f'EPSG:{dst_epsg}')), # destination coordinate system. Note: not using dst_crs as pyproj will complain 
                                # with future warning about soon to be unsupported authority/code string
        always_xy=True #True keeps the lon first in the coordinate
        ) 

    # Transform all points in the polygon from src_crs to dst_crs
    transformed_shape = transform(project.transform, shape)  # apply projection

    return transformed_shape 