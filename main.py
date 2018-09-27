import sys

def dfs_tree(root, get_children):
    """
    Run dfs over the tree rooted at r. Return a tuple (n_nodes, n_leaves) where
    n_nodes is the number of nodes in the tree rooted at r and n_leaves is the
    number of leaves in the tree rooted at r.

    The graph rooted at r MUST be a tree, otherwise dfs_tree will endless loop.
    """

    children = get_children(root)
    if len(children) == 0: return 1, 1

    n_nodes, n_leaves = 0, 0
    for x,y in [dfs_tree(child, get_children) for child in children]:
        n_nodes  += x
        n_leaves += y

    return n_nodes + 1, n_leaves


def solve(g_root, n_rows, n_cols):
    """
    start: A node in G
    n_rows: Number of rows in G
    n_cols: Number of cols in G

    Given the graph G, an n_rows by n_cols grid of nodes, return the tuple (x,
    y) where x is the number of non intersecting paths rooted at g_root and y
    is the number of terminal non intersecting paths rooted at g_root.

    Nodes in G are numbered 0 to n_rows*n_cols - 1, starting from the top left and
    continuing down like so
    0  1  2  3  4
    5  6  7  8  9
    10 11 12 13 14

    To count paths, consider the graph G' where each node in G' represents a
    non-intersecting path in G. Edges in G' work as expected. A node x in G' is
    the child of a node y in G' if adding a node in G to the path y gives you
    x. Each leaf node in G' represents a terminal non-intersecting path in G.
    Return the number of nodes and leaf nodes in G'

    To remain similar to the go code, represent a node in G' by the tuple
    (x, y), where x and y are integers. x is the last node in the path and y is
    a bitmask which is the set of nodes that have been visited along this path.

    NOTE: Even though G' is a tree, a node might appear twice in G'. For
    example if G is a 3x3 grid, there are 2 ways to go from the bottom left
    corner to the top right. This is a potential source for optimization.
    """

    def get_children_gprime( (n, visited_bit_mask) ):
        """
        Given a node in G', return it's children.
        """
        r, c = n // n_cols, n % n_cols

        children = ( (r + 1, c), (r - 1, c), (r, c + 1), (r, c - 1) )
        children = ( (r, c) for (r, c) in children if 0 <= r < n_rows and 0 <= c < n_cols )
        children = ( r * n_cols + c for (r, c) in children )
        children = ( n for n in children if  (1 << n) & visited_bit_mask == 0 )
        children = [ (n, visited_bit_mask | (1 << n)) for n in children ]
        return children

    g_prime_root = (g_root, 1 << g_root)
    return dfs_tree(g_prime_root, get_children_gprime)



if __name__ == '__main__':
    n_rows, n_cols = int(sys.argv[1]), int(sys.argv[2])
    start = 0 if len(sys.argv) < 4 else int(sys.argv[3])
    self_avoiding_paths, terminal_self_avoiding_paths = solve(start, n_rows, n_cols)
    print('self_avoiding_paths {}\tterminal_self_avoiding_paths {}'.format(
              self_avoiding_paths, terminal_self_avoiding_paths
          ))
