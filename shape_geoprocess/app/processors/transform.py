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