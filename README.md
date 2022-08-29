## What are the functions?
* file（文件操作）
    * func
        * GetFiles （读取指定目录下的文件）
        * CopyFolder （复制整个文件夹）
        * CopyFile （复制文件）
        * ClearFile （清空目录下的所有文件）
        * IsExists （判断路径是否存在）
        * Delete （删除文件）
        * WriteString （写入文件）
        * AppendString （追加文件）
        * AppendLine （换行追加文件）
        * AppendAllLine （换行追加文件）
        * CreateDir766 （创建所有目录，权限为766）
        * CreateDir （创建所有目录）
        * ReadString （读文件内容）
        * ReadAllLines （读文件内容，按行返回数组）
* encrypt（加密操作）
    * Md5 （对字符串做MD5加密）
* exec（shell）
    * RunShell （执行shell命令）
    * RunShellContext （执行shell命令）
* str（字符串操作）
    * CutRight （裁剪末尾标签）
    * MapToStringList （将map转成字符串数组）
    * ToDateTime （将时间转换为yyyy-MM-dd HH:mm:ss）
* http（http操作）
    * Post （http post，支持超时设置）
    * PostForm （http post，默认x-www-form-urlencoded）
    * PostFormWithoutBody （http post，默认x-www-form-urlencoded）
    * PostJson （Post方式将结果反序列化成TReturn）
    * Get （http get，支持超时设置）
    * GetForm （http get，默认x-www-form-urlencoded）
    * GetFormWithoutBody （http get，默认x-www-form-urlencoded）
    * GetJson （Get方式将结果反序列化成TReturn）
    * AddHttpPrefix（添加http前缀）
    * AddHttpsPrefix（添加https前缀）
* times（时间操作）
    * GetTime（根据time.Duration转换成天、小时、分钟、秒）
    * GetDesc（返回时间中文的描述）
    * GetSubDesc（返回时间中文的描述）
    * GetDate（获取当前日期）
* snowflake（雪花算法）
  * Init（全局初始化一次)
  * GenerateId（生成唯一ID）