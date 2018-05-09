package optimizer

import (
	"sqlparser"
	"log"
)
type optimizer struct {
	stmt     sqlparser.Statement
}
func Optimizer(stmt sqlparser.Statement)  {
	opt := newOptimizer(stmt)
	_ = sqlparser.Walk(opt.WalkStatement, stmt)
}
func newOptimizer(stmt sqlparser.Statement) *optimizer {
	return &optimizer{stmt:stmt};
}
func fetchColNamefromAliasExpr(node sqlparser.SQLNode)(colname *sqlparser.ColName){

	sqlparser.Walk(func(alias_node sqlparser.SQLNode) (kontinue bool, err error) {

		switch alias_node:= alias_node.(type) {
		case *sqlparser.ColName:
			colname = alias_node
			return false,nil
		}
		return true,nil
	}, node)
	return colname;
}
func fetchTableName(node sqlparser.SQLNode)(tabname *sqlparser.TableName)  {
	sqlparser.Walk(func(table_node sqlparser.SQLNode) (kontinue bool, err error) {

		switch table_node:= table_node.(type) {
		case sqlparser.TableName:
			tabname = &table_node
			return false,nil
		}
		return true,nil
	}, node)
	return tabname
}
func (opt *optimizer)WalkStatement(node sqlparser.SQLNode) (bool, error)  {
	switch node := node.(type) {
	case *sqlparser.Select:
		//group by 没有order by ,则添加 order by null
		if node.GroupBy != nil && node.OrderBy ==nil {
			//node.OrderBy = &sqlparser.OrderBy{}
			order := sqlparser.Order{Expr:&sqlparser.NullVal{}}
			node.AddOrder(&order)
		}
		// TODO 改写in语句为exists语句
		if len(node.From) == 1 {
			if tabexpr, ok:= (node.From[0]).(*sqlparser.AliasedTableExpr);ok{

				if simpleTbaleExpr,ok:=tabexpr.Expr.(sqlparser.SimpleTableExpr); ok{
					if tableName,ok := simpleTbaleExpr.(sqlparser.TableName); ok{

						sqlparser.Walk(func(node_ sqlparser.SQLNode) (kontinue bool, err error) {
							switch node_one_expr:=node_.(type){
							case sqlparser.Expr:
								if node_one,ok:=node_one_expr.(*sqlparser.ComparisonExpr); ok {
									if  node_left, yes := node_one.Left.(*sqlparser.ColName);yes &&node_one.Operator == sqlparser.InStr {
										if subquery, ok := node_one.Right.(*sqlparser.Subquery); ok {

											AliasExpr := sqlparser.AliasedExpr{Expr: sqlparser.NewIntVal([]byte("1")), As: sqlparser.ColIdent{}}
											// subquery.Select
											if subselect, ok:=subquery.Select.(*sqlparser.Select); ok{
												//  从子查询中读取列
												col_name := fetchColNamefromAliasExpr(subselect.SelectExprs[0]);

												subquery_tablename := fetchTableName(subselect.From[0])
												col_name.Qualifier = *subquery_tablename

												subselect.SelectExprs = sqlparser.SelectExprs{&AliasExpr}
												node_left.Qualifier = tableName
												compareExpr := sqlparser.ComparisonExpr{Operator:sqlparser.EqualStr, Left:node_left,Right:col_name}
												subselect.AddWhere(&compareExpr)
												node_one_expr = &sqlparser.ExistsExpr{Subquery:subquery}
												if (node_one == node.Where.Expr){
													node.Where = sqlparser.NewWhere(sqlparser.WhereStr, node_one_expr)
												}


											}

											return false, nil
										}
									}
								}

							}
							return true, nil
						}, node);
						log.Println(node)
					}
				}

			}
		}

		_ = sqlparser.Walk(opt.WalkOpt, node)
	case *sqlparser.Delete,*sqlparser.Update:

		_ = sqlparser.Walk(opt.WalkOpt, node)

	// Don't continue
	return false, nil

	}
	return true, nil
}
func check_has_colname(node sqlparser.Expr) bool{
	has := false
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type){
		case *sqlparser.ColName:
			has = true;
			return false, nil
		}
		return true, nil
	}, node);
	return has
}
/**

 */
func (opt *optimizer) optCompExpr(node *sqlparser.ComparisonExpr){

	right := node.Right
	switch right :=right.(type) {
		case *sqlparser.ParenExpr:

			if _, ok := right.Expr.(*sqlparser.ComparisonExpr);!ok{
				node.Right = right.Expr;
			}
		case *sqlparser.BinaryExpr:
			// 如果右边有colname,左右互换
			if(check_has_colname(right)){
				 node.Right = node.Left
				 node.Left = right
			}
	}
	switch left := (node.Left).(type){
	case *sqlparser.BinaryExpr:
		// 检查colname是否出现两次,出现两次不优化,不是+-*/  不优化
		count_colname := 0;
		has_other:=false
		binop_flag := true;
		sqlparser.Walk(func(node_left sqlparser.SQLNode) (kontinue bool, err error) {
			switch node_left:=node_left.(type){
			case *sqlparser.ColName:
				count_colname++;
				//子项目不在递归
				return false, nil
			case *sqlparser.SQLVal:
			case *sqlparser.ParenExpr:
			case *sqlparser.BinaryExpr:{
				switch node_left.Operator{
				case sqlparser.MinusStr:
				case sqlparser.PlusStr:
				case sqlparser.MultStr:
				case sqlparser.DivStr:
					binop_flag = true
				default:
					binop_flag = false
				}
			}
			default:
				has_other = true;
			}
			return true, nil
		}, left);
		if (count_colname == 1 && !has_other && binop_flag){
			go_on_flag:=true;
			for go_on_flag {
				root := left;
				op:=sqlparser.PlusStr;
				if (root.Operator == sqlparser.PlusStr){
					op = sqlparser.MinusStr;
				}
				if (root.Operator == sqlparser.MinusStr){
					op = sqlparser.PlusStr;
				}
				if (root.Operator == sqlparser.MultStr){
					op = sqlparser.DivStr;
				}
				if (root.Operator == sqlparser.DivStr){
					op = sqlparser.MultStr;
				}

				if (check_has_colname(root.Left)){

					switch one_node:=root.Left.(type){
					case *sqlparser.BinaryExpr:
						node.Right = &sqlparser.BinaryExpr{Left:&sqlparser.ParenExpr{node.Right}, Operator:op, Right:root.Right}
						left = one_node;
					case *sqlparser.ColName:
						//left.Right = &sqlparser.BinaryExpr{Left:left.Right, Operator:op, Right:one_node.Left}
						node.Right = &sqlparser.BinaryExpr{Left:&sqlparser.ParenExpr{node.Right}, Operator:op, Right:root.Right}
						node.Left = one_node;
						go_on_flag = false;

					}
				}else if (check_has_colname(root.Right)){
					switch one_node:=root.Right.(type){
					case *sqlparser.BinaryExpr:
						node.Right = &sqlparser.BinaryExpr{Left:&sqlparser.ParenExpr{node.Right}, Operator:op, Right:root.Left}
						left = one_node;
					case *sqlparser.ColName:
						//left.Right = &sqlparser.BinaryExpr{Left:left.Right, Operator:op, Right:one_node.Left}
						node.Right = &sqlparser.BinaryExpr{Left:&sqlparser.ParenExpr{node.Right}, Operator:op, Right:root.Left}
						node.Left = one_node;
						go_on_flag = false;
					}
				}

			}

		}
	}
}
func (opt *optimizer) WalkOpt(node sqlparser.SQLNode) (bool, error) {
	switch node := node.(type) {
	case *sqlparser.SQLVal:
	case *sqlparser.ComparisonExpr:
		opt.optCompExpr(node)
	}
	return true, nil
}