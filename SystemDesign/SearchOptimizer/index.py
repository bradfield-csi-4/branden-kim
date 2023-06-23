import csv
import pprint

inverted_index = {}
pp = pprint.PrettyPrinter(indent=4)

# Python function to print permutations of a given list
def permutation(lst):
 
    # If lst is empty then there are no permutations
    if len(lst) == 0:
        return []
 
    # If there is only one element in lst then, only
    # one permutation is possible
    if len(lst) == 1:
        return [lst]
 
    # Find the permutations for lst if there are
    # more than 1 characters
 
    l = [] # empty list that will store current permutation
 
    # Iterate the input(lst) and calculate the permutation
    for i in range(len(lst)):
       m = lst[i]
 
       # Extract lst[i] or m from the list.  remLst is
       # remaining list
       remLst = lst[:i] + lst[i+1:]
 
       # Generating all permutations where m is first
       # element
       for p in permutation(remLst):
           l.append([m] + p)
    return l

def main():
    with open("GNAF_CORE.psv") as csvfile:
        reader = csv.reader(csvfile, delimiter="|")

        count = 0
        for row in reader:
            if count == 0:
                print("COLUMNS: ")
            
            key = row[12:18]
            key = " ".join(key).split()
            value = row[25:]

            street_name = row[12]
            street_type = row[13]
            street_suffix = row[14]
            locality_name = row[15]
            state = row[16]
            postcode = row[17]
            
            longitude = row[25]
            latitude = row[26]
            
            if count != 0:
#                print_format(key, value)
#                inverted_index[" ".join(key)] = [tuple(value)] 
#                print(key)
                for p in permutation(key):
#                    print(p)
                    inverted_index[" ".join(p)] = [tuple(value)] 
            else:
                print(" ".join(key) + " " + " ".join(value))

            count+=1
            if count == 10:
                break

    pp.pprint(inverted_index)

def print_format(key, value):
    print("KEY:" + " ".join(key) + " VALUE: " + "(" + ",".join(value) + ")")


if __name__ == "__main__":
    main()
