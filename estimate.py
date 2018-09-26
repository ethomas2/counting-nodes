import random
from itertools import *

"""
Take a bunch of random walks through G. Keep track of the number of children at
each node. Keep track of the branching factor at each point and the depth of
each walk. Continuously output the average branching factor (across all nodes)
and average depth across all walks. The estimate is average_branching_factor **
average_depth.
"""
N, M = 4, 8
start = (0, 0)

def get_options( (x, y), visited):
    return [(xprime, yprime) for (xprime, yprime) in [(x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)]
            if 0 <= xprime < N and 0 <= yprime < M and (xprime, yprime) not in visited]

def gen_random_path():
    node = start
    visited = set(node)
    path = [node]
    branching_factors = []

    while True:
        opts = get_options(node, visited)
        if len(opts) == 0: return branching_factors, path
        branching_factors.append( len(opts) )
        node = random.choice(opts)
        path.append(node)
        visited.add(node)

branching_factors_sum = 0
branching_factors_count = 0
depths_sum = 0
depths_total = 0

for i in count():
    branching_factors, path = gen_random_path()
    branching_factors_sum += sum(branching_factors)
    branching_factors_count += len(branching_factors)
    depths_sum += len(path)
    depths_total += 1

    if i % 10000 == 0:
        bf_avg, depth_avg = float(branching_factors_sum) / branching_factors_count , float(depths_sum) / depths_total
        estimate = bf_avg ** depth_avg
        print('Round {i}\tbf_avg {bf_avg}\tdepth_avg {depth_avg}\test {estimate}'.format(
                    i=i, bf_avg=bf_avg, depth_avg=depth_avg, estimate=estimate,
              ))
