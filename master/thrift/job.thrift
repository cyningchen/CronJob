namespace go rpc
//1、定义数据结构
struct DataJob{
    1: string Name
    2: string Command
    3: string Expr
}

struct Result{
    1: list<DataJob> Job
}

//2、定义服务
service JobService{
    Result ListJob()
}

//3、 定义服务方法

