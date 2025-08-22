import pytest
from grpclib.client import Channel
import db_gateway_pb2 as pb
import db_gateway_grpc as pb_grpc
import uuid
import asyncio


@pytest.mark.asyncio
async def test_ping():
    """
    测试 Ping GRPC 方法
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)
        response = await stub.Ping(pb.PingRequest())
        assert response is not None


@pytest.mark.asyncio
async def test_mysql_create_and_drop_table():
    """
    测试 MySQL 创建和删除表
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        # 生成随机表名
        table_name = f"test_table_{uuid.uuid4().hex}"

        # 创建表的 SQL 请求
        create_table_sql = f"CREATE TABLE {table_name} (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))"
        create_request = pb.SqlRequest(
            sql=create_table_sql,
            db=pb.SqlDatabases.Masterdb,
            commit=True,
            get_row_count=False,
            get_last_insert_id=False
        )
        create_response = await stub.Mysql(create_request)
        assert create_response.result.code == 800, f"创建表失败: {create_response.result.msg}"
        print(f"表 {table_name} 创建成功")

        # 删除表的 SQL 请求
        drop_table_sql = f"DROP TABLE {table_name}"
        drop_request = pb.SqlRequest(
            sql=drop_table_sql,
            db=pb.SqlDatabases.Masterdb,
            commit=True,
            get_row_count=False,
            get_last_insert_id=False
        )
        drop_response = await stub.Mysql(drop_request)
        assert drop_response.result.code == 800, f"删除表失败: {drop_response.result.msg}"
        print(f"表 {table_name} 删除成功")


@pytest.mark.asyncio
async def test_redis_set_get_del_string():
    """
    测试 Redis 字符串的设置、获取和删除
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        key = f"test_key_{uuid.uuid4().hex}"
        value = "test_value"
        dbid = 0
        ttl = 60

        # 设置字符串
        set_request = pb.RedisSetStringRequest(
            DBID=dbid, key=key, value=value, ttl=ttl)
        set_response = await stub.RedisSet(set_request)
        assert set_response.result.code == 800, f"Redis 设置字符串失败: {set_response.result.msg}"
        print(f"Redis 键 {key} 设置成功")

        # 获取字符串
        get_request = pb.RedisGetStringRequest(DBID=dbid, key=key)
        get_response = await stub.RedisGet(get_request)
        assert get_response.result.code == 800, f"Redis 获取字符串失败: {get_response.result.msg}"
        assert get_response.value == value, "获取到的值与设置的值不匹配"
        print(f"Redis 键 {key} 获取成功，值为 {get_response.value}")

        # 删除字符串
        del_request = pb.RedisDelRequest(DBID=dbid, key=key)
        del_response = await stub.RedisDel(del_request)
        assert del_response.result.code == 800, f"Redis 删除字符串失败: {del_response.result.msg}"
        print(f"Redis 键 {key} 删除成功")


@pytest.mark.asyncio
async def test_redis_set_get_del_bytes():
    """
    测试 Redis 字节数组的设置、获取和删除
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        key = f"test_bytes_key_{uuid.uuid4().hex}"
        value = b"test_bytes_value"
        dbid = 0
        ttl = 60

        # 设置字节数组
        set_request = pb.RedisSetBytesRequest(
            DBID=dbid, key=key, value=value, ttl=ttl)
        set_response = await stub.RedisBSet(set_request)
        assert set_response.result.code == 800, f"Redis 设置字节数组失败: {set_response.result.msg}"
        print(f"Redis 字节数组键 {key} 设置成功")

        # 获取字节数组
        get_request = pb.RedisGetBytesRequest(DBID=dbid, key=key)
        get_response = await stub.RedisBGet(get_request)
        assert get_response.result.code == 800, f"Redis 获取字节数组失败: {get_response.result.msg}"
        assert get_response.value == value, "获取到的字节数组与设置的值不匹配"
        print(f"Redis 字节数组键 {key} 获取成功，值为 {get_response.value}")

        # 删除字节数组
        del_request = pb.RedisDelRequest(DBID=dbid, key=key)
        del_response = await stub.RedisDel(del_request)
        assert del_response.result.code == 800, f"Redis 删除字节数组失败: {del_response.result.msg}"
        print(f"Redis 字节数组键 {key} 删除成功")


@pytest.mark.asyncio
async def test_mysql_insert_select_update_delete():
    """
    测试 MySQL 插入、查询、更新和删除数据
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        table_name = f"test_data_table_{uuid.uuid4().hex}"
        create_table_sql = f"CREATE TABLE {table_name} (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), age INT)"
        create_request = pb.SqlRequest(
            sql=create_table_sql,
            db=pb.SqlDatabases.Masterdb,
            commit=True
        )
        create_response = await stub.Mysql(create_request)
        assert create_response.result.code == 800, f"创建表失败: {create_response.result.msg}"
        print(f"表 {table_name} 创建成功")

        try:
            # 插入数据
            insert_sql = f"INSERT INTO {table_name} (name, age) VALUES (?, ?)"
            insert_request = pb.SqlRequest(
                sql=insert_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(str="Alice"),
                    pb.InterFaceType(int32=30)
                ],
                commit=True,
                get_last_insert_id=True,
                get_row_count=True
            )
            insert_response = await stub.Mysql(insert_request)
            assert insert_response.result.code == 800, f"插入数据失败: {insert_response.result.msg}"
            assert insert_response.rows_affected == 1, "插入数据影响行数不正确"
            assert insert_response.last_insert_id > 0, "未获取到插入ID"
            inserted_id = insert_response.last_insert_id
            print(f"数据插入成功，ID: {inserted_id}")

            # 查询数据
            select_sql = f"SELECT id, name, age FROM {table_name} WHERE id = ?"
            select_request = pb.SqlRequest(
                sql=select_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(int64=inserted_id)
                ]
            )
            select_response = await stub.Mysql(select_request)
            assert select_response.result.code == 800, f"查询数据失败: {select_response.result.msg}"
            assert len(select_response.data) == 1, "查询结果行数不正确"
            assert select_response.data[0].result[0].int32 == inserted_id
            assert select_response.data[0].result[1].str == "Alice"
            assert select_response.data[0].result[2].int32 == 30
            print(f"数据查询成功: {select_response.data[0].result}")

            # 更新数据
            update_sql = f"UPDATE {table_name} SET age = ? WHERE id = ?"
            update_request = pb.SqlRequest(
                sql=update_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(int32=31),
                    pb.InterFaceType(int64=inserted_id)
                ],
                commit=True,
                get_row_count=True
            )
            update_response = await stub.Mysql(update_request)
            assert update_response.result.code == 800, f"更新数据失败: {update_response.result.msg}"
            assert update_response.rows_affected == 1, "更新数据影响行数不正确"
            print(f"数据更新成功，ID: {inserted_id}")

            # 再次查询验证更新
            select_response_after_update = await stub.Mysql(select_request)
            assert select_response_after_update.result.code == 800, f"更新后查询数据失败: {select_response_after_update.result.msg}"
            assert len(select_response_after_update.data) == 1
            assert select_response_after_update.data[0].result[2].int32 == 31
            print(f"更新后数据验证成功: {select_response_after_update.data[0].result}")

            # 删除数据
            delete_sql = f"DELETE FROM {table_name} WHERE id = ?"
            delete_request = pb.SqlRequest(
                sql=delete_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(int64=inserted_id)
                ],
                commit=True,
                get_row_count=True
            )
            delete_response = await stub.Mysql(delete_request)
            assert delete_response.result.code == 800, f"删除数据失败: {delete_response.result.msg}"
            assert delete_response.rows_affected == 1, "删除数据影响行数不正确"
            print(f"数据删除成功，ID: {inserted_id}")

            # 再次查询验证删除
            select_response_after_delete = await stub.Mysql(select_request)
            assert select_response_after_delete.result.code == 800, f"删除后查询数据失败: {select_response_after_delete.result.msg}"
            assert len(select_response_after_delete.data) == 0, "数据未被完全删除"
            print(f"删除后数据验证成功，无数据")

        finally:
            # 清理：删除表
            drop_table_sql = f"DROP TABLE {table_name}"
            drop_request = pb.SqlRequest(
                sql=drop_table_sql,
                db=pb.SqlDatabases.Masterdb,
                commit=True
            )
            drop_response = await stub.Mysql(drop_request)
            assert drop_response.result.code == 800, f"删除表失败: {drop_response.result.msg}"
            print(f"表 {table_name} 清理成功")


@pytest.mark.asyncio
async def test_mysql_datatypes():
    """
    测试 MySQL 不同数据类型的插入和查询
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        table_name = f"test_datatypes_table_{uuid.uuid4().hex}"
        create_table_sql = f"""CREATE TABLE {table_name} (
            id INT AUTO_INCREMENT PRIMARY KEY,
            str_col VARCHAR(255),
            int32_col INT,
            int64_col BIGINT,
            bool_col BOOLEAN,
            float_col FLOAT,
            double_col DOUBLE,
            blob_col BLOB
        )"""
        create_request = pb.SqlRequest(
            sql=create_table_sql,
            db=pb.SqlDatabases.Masterdb,
            commit=True
        )
        create_response = await stub.Mysql(create_request)
        assert create_response.result.code == 800, f"创建表失败: {create_response.result.msg}"
        print(f"表 {table_name} 创建成功")

        try:
            # 插入不同数据类型的数据
            insert_sql = f"""INSERT INTO {table_name} (
                str_col, int32_col, int64_col, bool_col, float_col, double_col, blob_col) VALUES (?, ?, ?, ?, ?, ?, ?)"""
            insert_request = pb.SqlRequest(
                sql=insert_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(str="test_string"),
                    pb.InterFaceType(int32=123),
                    pb.InterFaceType(int64=456789012345),
                    pb.InterFaceType(bool=True),
                    pb.InterFaceType(float=1.23),
                    pb.InterFaceType(double=4.56789),
                    pb.InterFaceType(blob=b"test_blob_data")
                ],
                commit=True,
                get_last_insert_id=True
            )
            insert_response = await stub.Mysql(insert_request)
            assert insert_response.result.code == 800, f"插入数据失败: {insert_response.result.msg}"
            assert insert_response.last_insert_id > 0, "未获取到插入ID"
            inserted_id = insert_response.last_insert_id
            print(f"不同数据类型数据插入成功，ID: {inserted_id}")

            # 查询数据并验证
            select_sql = f"SELECT * FROM {table_name} WHERE id = ?"
            select_request = pb.SqlRequest(
                sql=select_sql,
                db=pb.SqlDatabases.Masterdb,
                params=[
                    pb.InterFaceType(int64=inserted_id)
                ]
            )
            select_response = await stub.Mysql(select_request)
            assert select_response.result.code == 800, f"查询数据失败: {select_response.result.msg}"
            assert len(select_response.data) == 1, "查询结果行数不正确"
            row = select_response.data[0].result

            print(row)

            assert row[1].str == "test_string"
            assert row[2].int32 == 123
            assert row[3].int64 == 456789012345
            assert row[4].int32 == 1
            assert abs(row[5].float - 1.23) < 0.0001  # 浮点数比较
            assert abs(row[6].double - 4.56789) < 0.0000001  # 双精度浮点数比较
            assert row[7].blob == b"test_blob_data"
            print(f"不同数据类型数据查询验证成功: {row}")

        finally:
            # 清理：删除表
            drop_table_sql = f"DROP TABLE {table_name}"
            drop_request = pb.SqlRequest(
                sql=drop_table_sql,
                db=pb.SqlDatabases.Masterdb,
                commit=True
            )
            drop_response = await stub.Mysql(drop_request)
            assert drop_response.result.code == 800, f"删除表失败: {drop_response.result.msg}"
            print(f"表 {table_name} 清理成功")


@pytest.mark.asyncio
async def test_redis_ttl():
    """
    测试 Redis 键的 TTL (Time To Live) 功能
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        key = f"test_ttl_key_{uuid.uuid4().hex}"
        value = "ttl_value"
        dbid = 0
        ttl = 2  # 设置一个短的 TTL (2 秒)

        # 设置带有 TTL 的字符串
        set_request = pb.RedisSetStringRequest(
            DBID=dbid, key=key, value=value, ttl=ttl)
        set_response = await stub.RedisSet(set_request)
        assert set_response.result.code == 800, f"Redis 设置带有 TTL 的字符串失败: {set_response.result.msg}"
        print(f"Redis 键 {key} 设置成功，TTL: {ttl} 秒")

        # 立即获取，应该能获取到值
        get_request_initial = pb.RedisGetStringRequest(DBID=dbid, key=key)
        get_response_initial = await stub.RedisGet(get_request_initial)
        assert get_response_initial.result.code == 800, f"Redis 获取字符串失败: {get_response_initial.result.msg}"
        assert get_response_initial.value == value, "立即获取到的值与设置的值不匹配"
        print(f"Redis 键 {key} 立即获取成功，值为 {get_response_initial.value}")

        # 等待 TTL 过期
        await asyncio.sleep(ttl + 1)  # 等待比 TTL 稍长的时间

        # 再次获取，应该获取不到值
        get_request_expired = pb.RedisGetStringRequest(DBID=dbid, key=key)
        get_response_expired = await stub.RedisGet(get_request_expired)
        # 假设获取不到值时，code可能不是800，或者value为空。这里需要根据实际服务行为调整断言
        # 如果服务返回特定错误码表示键不存在，则断言该错误码
        # 如果服务返回成功但value为空，则断言value为空
        # 暂时断言value为空，如果实际行为不同，需要调整
        assert get_response_expired.value == "", "过期键仍然获取到值"
        print(f"Redis 键 {key} 过期后获取成功，值为 {get_response_expired.value} (预期为空)")

        # 清理 (如果键已经过期，删除操作可能仍然返回成功，或者返回键不存在的错误码)
        del_request = pb.RedisDelRequest(DBID=dbid, key=key)
        del_response = await stub.RedisDel(del_request)
        # 这里根据实际服务行为调整断言，如果键已过期，删除可能返回成功或键不存在的错误码
        assert del_response.result.code == 800, f"Redis 删除过期键失败: {del_response.result.msg}"
        print(f"Redis 键 {key} 清理成功")


@pytest.mark.asyncio
async def test_mysql_invalid_sql():
    """
    测试 MySQL 执行无效 SQL 语句
    """
    async with Channel('127.0.0.1', 50051) as channel:
        stub = pb_grpc.StealthIMDBGatewayStub(channel)

        invalid_sql = "SELECT * FROM non_existent_table"
        request = pb.SqlRequest(
            sql=invalid_sql,
            db=pb.SqlDatabases.Masterdb,
            commit=False
        )
        response = await stub.Mysql(request)

        # 预期会返回 DBGatewaySQLExecuteError (2002)
        assert response.result.code == 2002, f"无效 SQL 语句未返回预期错误码: {response.result.msg}"
        print(
            f"执行无效 SQL 语句成功捕获错误。Code: {response.result.code}, Msg: {response.result.msg}")
