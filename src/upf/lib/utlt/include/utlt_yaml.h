#ifndef __UTLT_YAML_H__
#define __UTLT_YAML_H__

#include <yaml.h>

#include "utlt_debug.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

#define GET_KEY 0
#define GET_VALUE 1

typedef struct _YamlIter {
    yaml_document_t     *document;
    yaml_node_t         *node;
    yaml_node_item_t    *seqItem;  /* An element of a sequence node. */
    yaml_node_pair_t    *mapPair;  /* An element of a mapping node. */
} YamlIter;

void YamlIterSetElement(YamlIter *iter, yaml_node_type_t type);
void YamlIterInit(YamlIter *iter, yaml_document_t *doc);
int YamlIterNext(YamlIter *iter);
int YamlIterType(YamlIter *iter);
void YamlIterChild(YamlIter *parent, YamlIter *child);
const char *YamlIterGet(YamlIter *iter, int getType);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif 