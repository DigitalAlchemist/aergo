/**
 * @file    check_stmt.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "check_id.h"
#include "check_exp.h"
#include "check_blk.h"

#include "check_stmt.h"

static int
stmt_if_check(check_t *check, ast_stmt_t *stmt)
{
    int i;
    ast_exp_t *cond_exp;
    meta_t *cond_meta;
    array_t *elif_stmts;

    ASSERT1(is_if_stmt(stmt), stmt->kind);
    ASSERT(stmt->u_if.cond_exp != NULL);

    cond_exp = stmt->u_if.cond_exp;
    cond_meta = &cond_exp->meta;

    CHECK(exp_check(check, cond_exp));

    if (!is_bool_meta(cond_meta))
        RETURN(ERROR_INVALID_COND_TYPE, &cond_exp->pos, meta_to_str(cond_meta));

    if (stmt->u_if.if_blk != NULL)
        blk_check(check, stmt->u_if.if_blk);

    elif_stmts = &stmt->u_if.elif_stmts;

    for (i = 0; i < array_size(elif_stmts); i++) {
        stmt_check(check, array_item(elif_stmts, i, ast_stmt_t));
    }

    if (stmt->u_if.else_blk != NULL)
        blk_check(check, stmt->u_if.else_blk);

    return NO_ERROR;
}

static int
stmt_loop_check_for(check_t *check, ast_stmt_t *stmt, char *begin_label,
                    char *end_label)
{
    ast_blk_t *blk = stmt->u_loop.blk;
    ast_exp_t *cond_exp;
    ast_exp_t *loop_exp;
    ast_stmt_t *null_stmt;

    if (stmt->u_loop.init_ids != NULL) {
        ASSERT(stmt->u_loop.init_exp == NULL);
        id_join_first(&blk->ids, stmt->u_loop.init_ids);
    }
    else {
        ast_exp_t *init_exp = stmt->u_loop.init_exp;

        if (init_exp != NULL) {
            ast_stmt_t *exp_stmt = stmt_new_exp(init_exp, &init_exp->pos);
            array_add_first(&blk->stmts, exp_stmt);
        }
    }

    cond_exp = stmt->u_loop.cond_exp;

    if (cond_exp != NULL) {
        ast_blk_t *if_blk;
        ast_stmt_t *goto_stmt;
        ast_exp_t *not_exp;
        ast_stmt_t *if_stmt;

        goto_stmt = stmt_new_goto(xstrdup(end_label), &cond_exp->pos);

        if_blk = blk_new_anon(&cond_exp->pos);
        array_add_last(&if_blk->stmts, goto_stmt);

        not_exp = exp_new_op(OP_NOT, cond_exp, NULL, &cond_exp->pos);

        if_stmt = stmt_new_if(not_exp, if_blk, &cond_exp->pos);
        array_add_first(&blk->stmts, if_stmt);
    }

    null_stmt = stmt_new_null(&stmt->pos);
    null_stmt->label = xstrdup(begin_label);

    array_add_first(&blk->stmts, null_stmt);

    loop_exp = stmt->u_loop.loop_exp;

    if (loop_exp != NULL) {
        ast_stmt_t *exp_stmt = stmt_new_exp(loop_exp, &loop_exp->pos);

        array_add_last(&blk->stmts, exp_stmt);
    }

    return NO_ERROR;
}

static int
stmt_loop_check_each(check_t *check, ast_stmt_t *stmt, char *begin_label,
                     char *end_label)
{
    char name[128];
    ast_id_t *id;
    ast_exp_t *inc_exp;
    ast_exp_t *arr_exp;
    ast_exp_t *assign_exp;
    ast_exp_t *loop_exp;
    ast_stmt_t *null_stmt;
    ast_blk_t *blk = stmt->u_loop.blk;
    src_pos_t *pos = &stmt->pos;

    loop_exp = stmt->u_loop.loop_exp;
    ASSERT(loop_exp != NULL);

    /* make "int i = 0" */
    snprintf(name, sizeof(name), "each_loop_idx_%d", blk->num);

    id = id_new_var(xstrdup(name), pos);

    id->u_var.type_exp = exp_new_type(TYPE_INT32, pos);
    id->u_var.size_exps = NULL;
    id->u_var.init_exp = exp_new_val(pos);
    value_set_int(&id->u_var.init_exp->u_val.val, 0);

    id_add_last(&blk->ids, id);

    inc_exp = exp_new_op(OP_INC, exp_new_id(xstrdup(name), pos), NULL, pos);
    arr_exp = exp_new_array(loop_exp, inc_exp, &loop_exp->pos);

    if (stmt->u_loop.init_ids != NULL) {
        int i;
        array_t *var_ids = stmt->u_loop.init_ids;

        /* make "variable = loop_exp[i++]" */
        for (i = 0; i < array_size(var_ids); i++) {
            ast_id_t *var_id = array_item(var_ids, i, ast_id_t);
            ast_exp_t *id_exp;

            id_exp = exp_new_id(var_id->name, pos);
            assign_exp =
                exp_new_op(OP_ASSIGN, id_exp, arr_exp, &loop_exp->pos);

            array_add_first(&blk->stmts, stmt_new_exp(assign_exp, pos));
        }

        id_join_first(&blk->ids, var_ids);
    }
    else {
        ast_exp_t *init_exp = stmt->u_loop.init_exp;

        ASSERT(init_exp != NULL);

        /* TODO: map iteration */
        if (is_tuple_exp(init_exp))
            RETURN(ERROR_NOT_SUPPORTED, &init_exp->pos);

        /* make "init_exp = loop_exp[i++]" */
        assign_exp =
            exp_new_op(OP_ASSIGN, init_exp, arr_exp, &loop_exp->pos);

        array_add_first(&blk->stmts, stmt_new_exp(assign_exp, pos));
    }

    null_stmt = stmt_new_null(&stmt->pos);
    null_stmt->label = xstrdup(begin_label);

    array_add_first(&blk->stmts, null_stmt);

    return NO_ERROR;
}

static int
stmt_loop_check(check_t *check, ast_stmt_t *stmt)
{
    char begin_label[128];
    char end_label[128];
    ast_stmt_t *goto_stmt;
    ast_stmt_t *null_stmt;
    ast_blk_t *blk;

    ASSERT1(is_loop_stmt(stmt), stmt->kind);

    if (stmt->u_loop.blk == NULL)
        stmt->u_loop.blk = blk_new_loop(&stmt->pos);

    blk = stmt->u_loop.blk;

    snprintf(begin_label, sizeof(begin_label), "for_loop_begin_%d", blk->num);
    snprintf(end_label, sizeof(end_label), "for_loop_end_%d", blk->num);

    switch (stmt->u_loop.kind) {
    case LOOP_FOR:
        stmt_loop_check_for(check, stmt, begin_label, end_label);
        break;

    case LOOP_EACH:
        stmt_loop_check_each(check, stmt, begin_label, end_label);
        break;

    default:
        ASSERT1(!"invalid loop", stmt->u_loop.kind);
    }

    goto_stmt = stmt_new_goto(xstrdup(begin_label), &stmt->pos);

    null_stmt = stmt_new_null(&stmt->pos);
    null_stmt->label = xstrdup(end_label);

    array_add_last(&blk->stmts, goto_stmt);
    array_add_last(&blk->stmts, null_stmt);

    blk_check(check, blk);

    return NO_ERROR;
}

static int
stmt_case_check(check_t *check, ast_stmt_t *stmt, meta_t *meta)
{
    int i;
    ast_exp_t *val_exp;
    array_t *stmts;

    ASSERT1(is_case_stmt(stmt), stmt->kind);

    val_exp = stmt->u_case.val_exp;

    if (val_exp != NULL) {
        meta_t *val_meta = &val_exp->meta;

        exp_check(check, val_exp);

        if (meta == NULL) {
            if (!is_bool_meta(val_meta))
                RETURN(ERROR_INVALID_COND_TYPE, &val_exp->pos,
                       meta_to_str(val_meta));
        }
        else if (!meta_equals(meta, val_meta)) {
            RETURN(ERROR_MISMATCHED_TYPE, &val_exp->pos,
                   meta_to_str(meta), meta_to_str(val_meta));
        }
    }

    stmts = stmt->u_case.stmts;

    for (i = 0; i < array_size(stmts); i++) {
        stmt_check(check, array_item(stmts, i, ast_stmt_t));
    }

    return NO_ERROR;
}

static int
stmt_switch_check(check_t *check, ast_stmt_t *stmt)
{
    int i;
    ast_exp_t *cond_exp;
    meta_t *cond_meta = NULL;
    array_t *case_stmts;

    ASSERT1(is_switch_stmt(stmt), stmt->kind);

    cond_exp = stmt->u_sw.cond_exp;

    if (cond_exp != NULL) {
        cond_meta = &cond_exp->meta;

        exp_check(check, cond_exp);

        if (!is_comparable_meta(cond_meta))
            RETURN(ERROR_NOT_COMPARABLE_TYPE, &cond_exp->pos,
                   meta_to_str(cond_meta));
    }

    case_stmts = stmt->u_sw.case_stmts;

    for (i = 0; i < array_size(case_stmts); i++) {
        stmt_case_check(check, array_item(case_stmts, i, ast_stmt_t),
                        cond_meta);
    }

    return NO_ERROR;
}

static int
stmt_return_check(check_t *check, ast_stmt_t *stmt)
{
    ast_id_t *fn_id;
    meta_t *fn_meta;
    ast_exp_t *arg_exp;

    ASSERT1(is_return_stmt(stmt), stmt->kind);
    ASSERT(check->fn_id != NULL);

    fn_id = check->fn_id;
    fn_meta = &fn_id->meta;

    ASSERT1(is_func_id(fn_id), fn_id->kind);

    arg_exp = stmt->u_ret.arg_exp;

    if (arg_exp != NULL) {
        ASSERT1(is_tuple_meta(fn_meta), fn_meta->type);

        if (is_void_meta(fn_meta))
            RETURN(ERROR_MISMATCHED_COUNT, &arg_exp->pos, 0,
                   meta_size(&arg_exp->meta));

        exp_check(check, arg_exp);

        if (is_tuple_meta(&arg_exp->meta)) {
            int i;
            array_t *arg_metas = arg_exp->meta.u_tup.metas;
            array_t *ret_metas = fn_meta->u_tup.metas;

            if (array_size(arg_metas) != array_size(ret_metas))
                RETURN(ERROR_MISMATCHED_COUNT, &arg_exp->pos,
                       array_size(ret_metas), array_size(arg_metas));

            for (i = 0; i < array_size(arg_metas); i++) {
                meta_t *arg_meta = array_item(arg_metas, i, meta_t);
                meta_t *ret_meta = array_item(ret_metas, i, meta_t);

                if (!meta_equals(ret_meta, arg_meta))
                    RETURN(ERROR_MISMATCHED_TYPE, &arg_exp->pos,
                           meta_to_str(ret_meta), meta_to_str(arg_meta));
            }
        }
        else {
            meta_t *arg_meta = &arg_exp->meta;
            array_t *ret_metas = fn_meta->u_tup.metas;
            meta_t *ret_meta;

            if (array_size(ret_metas) != 1)
                RETURN(ERROR_MISMATCHED_COUNT, &arg_exp->pos,
                       array_size(ret_metas), 1);

            ret_meta = array_item(fn_meta->u_tup.metas, 0, meta_t);

            if (!meta_equals(arg_meta, ret_meta))
                RETURN(ERROR_MISMATCHED_TYPE, &arg_exp->pos,
                       meta_to_str(ret_meta), meta_to_str(arg_meta));
        }
    }
    else if (!is_void_meta(fn_meta)) {
        RETURN(ERROR_MISMATCHED_COUNT, &stmt->pos, meta_size(fn_meta), 0);
    }

    return NO_ERROR;
}

static int
stmt_goto_check(check_t *check, ast_stmt_t *stmt)
{
    int i;
    int stmt_cnt;
    bool has_found = false;
    ast_blk_t *blk = check->blk;

    ASSERT1(is_goto_stmt(stmt), stmt->kind);
    ASSERT(stmt->u_goto.label != NULL);
    ASSERT(blk != NULL);

    do {
        stmt_cnt = array_size(&blk->stmts);

        for (i = 0; i < stmt_cnt; i++) {
            ast_stmt_t *prev = array_item(&blk->stmts, i, ast_stmt_t);

            if (prev->label != NULL &&
                strcmp(prev->label, stmt->u_goto.label) == 0) {
                has_found = true;
                break;
            }
        }
    } while ((blk = blk->up) != NULL);

    if (!has_found)
        RETURN(ERROR_UNDEFINED_LABEL, &stmt->pos, stmt->u_goto.label);

    return NO_ERROR;
}

static int
stmt_ddl_check(check_t *check, ast_stmt_t *stmt)
{
    ASSERT1(is_ddl_stmt(stmt), stmt->kind);
    ASSERT(stmt->u_ddl.ddl != NULL);

    return NO_ERROR;
}

static int
stmt_blk_check(check_t *check, ast_stmt_t *stmt)
{
    ASSERT1(is_blk_stmt(stmt), stmt->kind);

    if (stmt->u_blk.blk != NULL)
        blk_check(check, stmt->u_blk.blk);

    return NO_ERROR;
}

void
stmt_check(check_t *check, ast_stmt_t *stmt)
{
    switch (stmt->kind) {
    case STMT_NULL:
    case STMT_CONTINUE:
    case STMT_BREAK:
        break;

    case STMT_EXP:
        exp_check(check, stmt->u_exp.exp);
        break;

    case STMT_IF:
        stmt_if_check(check, stmt);
        break;

    case STMT_LOOP:
        stmt_loop_check(check, stmt);
        break;

    case STMT_SWITCH:
        stmt_switch_check(check, stmt);
        break;

    case STMT_RETURN:
        stmt_return_check(check, stmt);
        break;

    case STMT_GOTO:
        stmt_goto_check(check, stmt);
        break;

    case STMT_DDL:
        stmt_ddl_check(check, stmt);
        break;

    case STMT_BLK:
        stmt_blk_check(check, stmt);
        break;

    default:
        ASSERT1(!"invalid statement", stmt->kind);
    }
}

/* end of check_stmt.c */
