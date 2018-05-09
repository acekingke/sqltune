package optimizer

import (
	"testing"
	"sqlparser"
	"fmt"
)

func TestExprOpt(t *testing.T) {
	sql := "select * from t2 where a = (1=22)"
	want := "select * from t2 where a = (1=22)"

	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sqlparser.String(stmt))
	if sql != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}

}
func TestExprBinayLeft(t *testing.T)  {
	sql := "select * from t2 where a*2 = 1+2"
	want := "select * from t2 where a = (1 + 2) / 2"
	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}
}
func TestExprBinayLeft2(t *testing.T)  {
	sql := "select * from t2 where a*2 + 1 = 1+2"
	want := "select * from t2 where a = ((1 + 2) - 1) / 2"
	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}
}
func TestExprBinayLeft3(t *testing.T)  {
	sql := "select * from t2 where 2*a + 1 = 1+2"
	want := "select * from t2 where a = ((1 + 2) - 1) / 2"
	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}
}
func TestExprBinayLeft4(t *testing.T)  {
	sql := "select * from t2 where 1+2 =  2*a + 1"
	want := "select * from t2 where a = ((1 + 2) - 1) / 2"
	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}
}

func TestDeleteExprOpt(t *testing.T) {
	sql := "delete from t2 where a = (1+22)"
	want := "delete from t2 where a = 1 + 22"

	stmt, err := sqlparser.Parse(sql)
	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}

}
func TestGroupbyOpt(t *testing.T) {
	sql := "select b from t2 where a*2 + 1 = 1+2 group by b "
	want := "select b from t2 where a = ((1 + 2) - 1) / 2 group by b order by null"

	stmt, err := sqlparser.Parse(sql)

	Optimizer(stmt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sqlparser.String(stmt))
	if sqlparser.String(stmt) != want{
		t.Errorf("Append: %s, want %s", sql, want)
	}

}
func TestInSubquery(t *testing.T){
	sql := "select b from t2 where a in (select a from t3) "
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		t.Error(err)
	}
	Optimizer(stmt)
	fmt.Println("------------------")
	fmt.Println(sqlparser.String(stmt))
}
func TestOnVal(t *testing.T)  {
	sql := "select 1 from t2 "
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		t.Error(err)
	}
	Optimizer(stmt)
}