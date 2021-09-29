'''
Dissolve uses Shapely's Unary Union
https://shapely.readthedocs.io/en/stable/manual.html#shapely.ops.unary_union

Returns a representation of the union of the given geometric objects.
Areas of overlapping Polygons will get merged. LineStrings will get fully 
dissolved and noded. Duplicate Points will get merged. Because the union merges 
the areas of overlapping Polygons it can be used in an attempt to fix invalid MultiPolygons.
'''
from shapely.ops import unary_union

def process(shapes):
    return unary_union(shapes) 