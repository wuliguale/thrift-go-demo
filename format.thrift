

//定义结构体
struct Data {
    1: string text
}

//定义一个service，2个方法
service format_data {
    Data do_format(1:Data data),
    Data hello(1:Data data),
}



