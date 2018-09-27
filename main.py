import random
import sys

N, M = int(sys.argv[1]), int(sys.argv[2])
start = (0, 0)

def get_options( (x, y), visited_nodes):
    return [(xprime, yprime) for (xprime, yprime) in [(x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)]
            if 0 <= xprime < N and 0 <= yprime < M and (xprime, yprime) not in visited_nodes]

def dfs(path_so_far=None, visited_nodes=None):
    """
    This is dfs-ing through the graph G', NOT G. You can think of each
    invocation of this function as visiting a node in G'. G' is a tree, where G
    is not. path_so_far represents a path in G, which is a node in G'. Given a
    node n' in G', this function returns the tuple (x, y), where x is the total
    number of nodes in the tree rooted at n' and y is the number of terminal
    nodes at the tree rooted at n'.
    """
    if path_so_far is None:   path_so_far = [ start ]
    if visited_nodes is None: visited_nodes = set( [start] )

    curr_node = path_so_far[-1]
    opts      = get_options(curr_node, visited_nodes)

    if len(opts) == 0: return 1, 1

    total_paths, total_terminal_paths = 0, 0
    for opt in opts:
        n_paths, n_terminal_paths = dfs(path_so_far + [opt], visited_nodes | set( [opt] ) )
        total_paths               += n_paths
        total_terminal_paths      += n_terminal_paths

    return total_paths + 1, total_terminal_paths

print( dfs() )
