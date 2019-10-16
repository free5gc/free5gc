#ifndef __HASH_H__
#define __HASH_H__

#include "utlt_debug.h"
#include "utlt_buff.h"
#include "utlt_time.h"

#ifdef __cplusplus
extern "C" {
#endif

//ave hash compute the length automatically.
#define HASH_KEY_STRING     (-1)

typedef struct _HashEntry HashEntry;
typedef struct _Hash Hash;              // for hash tables.
typedef struct _HashIndex HashIndex;    // for scanning hash tables.
// Callback functions for calculating hash values.
typedef unsigned int (*HashFunc)(const char *key, int *klen);

// The internal form of a hash table.
struct _HashEntry {
    HashEntry       *next;
    unsigned int    hash;
    const void      *key;
    int             klen;
    const void      *val;
};

// Structure for iterating through a hash table.
struct _HashIndex {
    Hash            *ht;
    HashEntry       *this, *next;
    unsigned int    index;
};

struct _Hash {
    HashEntry       **array;
    HashIndex       iterator;  /* For HashFirst(NULL, ...) */
    unsigned int    count, max, seed;
    HashFunc        hashFunc;
    HashEntry       *entry;  /* List of hash entries of the hash table */
};

// The default hash function.
unsigned int CoreHashfuncDefault(const char *key, int *klen);

Hash *HashMake();
Hash *HashMakeCustom(HashFunc hashFunc);   // Create a hash table with a custom hash function
void HashDestroy(Hash *ht);
void HashSet(Hash *ht, const void *key, int klen, const void *val);
void *HashGet(Hash *ht, const void *key, int klen);

/**
 * Look up the value associated with a key in a hash table, or if none exists
 * associate a value.
 */
void *HashGetOrSet(Hash *ht, const void *key, int klen, const void *val);

// Iterating over the entries in a hash table.
HashIndex *HashFirst(Hash *ht);
HashIndex *HashNext(HashIndex *hi);

// Get the current entry's details from the iteration state.
void HashThis(HashIndex *hi, const void **key, int *klen, void **val);
const void *HashThisKey(HashIndex *hi);
int HashThisKeyLen(HashIndex *hi);
void *HashThisVal(HashIndex *hi);

unsigned int HashCount(Hash *ht);
void HashClear(Hash *ht);

// Declaration prototype for the iterator callback function of HashDo().
typedef int (HashDoCallbackFunc)(
        void *rec, const void *key, int klen, const void *value);

/** 
 * Iterate over a hash table running the provided function once for every
 * element in the hash table. The @param comp function will be invoked for
 * every element in the hash table.
 */
int HashDo(HashDoCallbackFunc *comp, void *rec, const Hash *ht);

#ifdef __cplusplus
}
#endif

#endif	/* __HASH_H__ */
