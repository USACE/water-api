'''
Simplify uses Shapely's simplify()

# https://shapely.readthedocs.io/en/stable/manual.html#object.simplify
All points in the simplified object will be within the tolerance distance of the 
original geometry. By default a slower algorithm is used that preserves topology. 
If preserve topology is set to False the much quicker Douglas-Peucker algorithm is used.
'''
import shapely

def process(shape):
    return shape.simplify(tolerance=0.5, preserve_topology=False) 