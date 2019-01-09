/**
 * @file    trans_blk.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "trans_id.h"
#include "trans_stmt.h"

#include "trans_blk.h"

void
blk_trans(trans_t *trans, ast_blk_t *blk)
{
    int i;
    int idx = 0;

    array_foreach(&blk->ids, i) {
        ast_id_t *id = array_get_id(&blk->ids, i);

        id_trans(trans, id);

        if (is_fn_id(id))
            id->idx = idx++;
    }

    array_foreach(&blk->stmts, i) {
        stmt_trans(trans, array_get_stmt(&blk->stmts, i));
    }

    /* Since there is no relation between general blocks and basic blocks,
     * we deliberately do not do anything about basic blocks here. */
}

/* end of trans_blk.c */