package main
import ("fmt"
	"sqlparser"
	"optimizer"
	"flag"
)


func main() {
	insql:=flag.String("insql", "select 1 from t", "insql=your input sql")
	flag.Parse();
	fmt.Println(*insql)
	stmt, err := sqlparser.Parse(*insql)

	if err != nil {
		// Do something with the err
		fmt.Println(err)
	}
	optimizer.Optimizer(stmt)
	fmt.Println("Output sql is:")
	fmt.Println(sqlparser.String(stmt))

}
