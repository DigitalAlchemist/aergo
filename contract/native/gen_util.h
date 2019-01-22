/**
 * @file    gen_util.h
 * @copyright defined in aergo/LICENSE.txt
 */

#ifndef _GEN_UTIL_H
#define _GEN_UTIL_H

#include "common.h"

#include "gen.h"
#include "binaryen-c.h"

#define i32_gen(gen, v)     BinaryenConst((gen)->module, BinaryenLiteralInt32(v))
#define i64_gen(gen, v)     BinaryenConst((gen)->module, BinaryenLiteralInt64(v))
#define f32_gen(gen, v)     BinaryenConst((gen)->module, BinaryenLiteralFloat32(v))
#define f64_gen(gen, v)     BinaryenConst((gen)->module, BinaryenLiteralFloat64(v))

#define meta_gen(meta)                                                                   \
    (is_array_meta(meta) ? BinaryenTypeInt32() : type_gen((meta)->type))

#ifndef _IR_SGMT_T
#define _IR_SGMT_T
typedef struct ir_sgmt_s ir_sgmt_t;
#endif /* ! _IR_SGMT_T */

void local_add(gen_t *gen, type_t type);
void instr_add(gen_t *gen, BinaryenExpressionRef instr);

BinaryenType type_gen(type_t type);

void table_gen(gen_t *gen, array_t *fns);
void sgmt_gen(gen_t *gen, ir_sgmt_t *sgmt);

#endif /* ! _GEN_UTIL_H */