package grpc

import (
	pb "StealthIMDB/StealthIM.DBGateway"
	"StealthIMDB/mysql"
	"context"
	"database/sql"
	"log"
	"time"
)

func isNil(i any) bool {
	if i == nil {
		return true
	}
	switch v := i.(type) {
	case *int, *float64, *bool, *string, *[]byte:
		return v == nil
	default:
		return false
	}
}

func (s *server) Mysql(ctx context.Context, in *pb.SqlRequest) (*pb.SqlResponse, error) {
	if cfg.GRPCProxy.Log {
		log.Printf("[GRPC][Mysql]Call \"%s\"\n", in.Sql)
	}
	var rowCount int64
	var lastInsertID int64

	db := mysql.GetConn(int32(in.Db) + 1)()
	if db == nil {
		return &pb.SqlResponse{Result: &pb.Result{Code: -1, Msg: "ConnectError"}}, nil
	}
	args := make([]any, len(in.Params))
	for i, obj := range in.Params {
		if x, ok := obj.Response.(*pb.InterFaceType_Int32); ok {
			args[i] = x.Int32
		} else if x, ok := obj.Response.(*pb.InterFaceType_Int64); ok {
			args[i] = x.Int64
		} else if x, ok := obj.Response.(*pb.InterFaceType_Str); ok {
			args[i] = x.Str
		} else if x, ok := obj.Response.(*pb.InterFaceType_Float); ok {
			args[i] = x.Float
		} else if x, ok := obj.Response.(*pb.InterFaceType_Double); ok {
			args[i] = x.Double
		} else if x, ok := obj.Response.(*pb.InterFaceType_Bool); ok {
			args[i] = x.Bool
		} else if x, ok := obj.Response.(*pb.InterFaceType_Blob); ok {
			args[i] = x.Blob
		} else {
			args[i] = nil
		}
	}
	if in.Commit {
		tx, err := db.Begin()
		if err != nil {
			return &pb.SqlResponse{Result: &pb.Result{Code: 1, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
		}
		var result sql.Result
		result, err = tx.Exec(in.Sql, args...)
		if err != nil {
			tx.Rollback()
			return &pb.SqlResponse{Result: &pb.Result{Code: -2, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
		}
		if in.GetRowCount {
			rowCount, err = result.RowsAffected()
			if err != nil {
				return &pb.SqlResponse{Result: &pb.Result{Code: -4, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
			}
		}
		if in.GetLastInsertId {
			lastInsertID, err = result.LastInsertId()
			if err != nil {
				return &pb.SqlResponse{Result: &pb.Result{Code: -5, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
			}
		}
		err = tx.Commit()
		if err != nil {
			return &pb.SqlResponse{Result: &pb.Result{Code: -3, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
		}
		return &pb.SqlResponse{Result: &pb.Result{Code: 0, Msg: ""}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
	}
	rows, err := db.Query(in.Sql, args...)
	if err != nil {
		return &pb.SqlResponse{Result: &pb.Result{Code: -1, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return &pb.SqlResponse{Result: &pb.Result{Code: -2, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return &pb.SqlResponse{Result: &pb.Result{Code: -3, Msg: err.Error()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
	}

	// 创建一个切片来保存每一行的数据
	datas := make([]*pb.SqlLine, 0)
	lncnt := 0
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		dataTmp := make([]*pb.InterFaceType, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// 读取当前行的数据
		err := rows.Scan(valuePtrs...)
		if err != nil {
			dataTmp = make([]*pb.InterFaceType, 1)
		}

		// 打印当前行的数据
		for i := range columns {
			var val any
			val = values[i]
			var isnull bool
			if isNil(val) {
				isnull = true
				dataTmp[i] = &pb.InterFaceType{Null: true}
				continue
			}
			// 根据列类型处理数据
			switch types[i].DatabaseTypeName() {
			case "VARCHAR", "CHAR", "TEXT":
				txt := ""
				b, ok := val.([]byte)
				if ok {
					txt = string(b)
				}
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Str{Str: txt}}
			case "BLOB":
				if isnull || val.([]byte) == nil {
					isnull = true
					val = make([]byte, 0)
				}
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Blob{Blob: val.([]byte)}, Null: isnull}
			case "TINYINT", "INT(8)", "SMALLINT", "INT(16)", "INT", "INTEGER", "MEDIUMINT", "INT(32)":
				if isnull {
					isnull = true
					val = make([]byte, 0)
				}
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Int32{Int32: int32(val.(int64))}, Null: isnull}
			case "BIGINT", "INT(64)":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Int64{Int64: val.(int64)}, Null: isnull}
			case "FLOAT":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Float{Float: val.(float32)}, Null: isnull}
			case "DOUBLE":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Double{Double: val.(float64)}, Null: isnull}
			case "DECIMAL":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Double{Double: val.(float64)}, Null: isnull}
			case "DATETIME", "TIMESTAMP":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Int64{Int64: val.(time.Time).Unix()}, Null: isnull}
			case "BOOLEAN", "TINYINT(1)":
				dataTmp[i] = &pb.InterFaceType{Response: &pb.InterFaceType_Bool{Bool: val.(bool)}, Null: isnull}
			default:
				return &pb.SqlResponse{Result: &pb.Result{Code: -4, Msg: "Unknown type: " + types[i].DatabaseTypeName()}, RowsAffected: rowCount, LastInsertId: lastInsertID}, nil
			}
		}
		datas = append(datas, &pb.SqlLine{Result: dataTmp})
		lncnt++
	}
	return &pb.SqlResponse{Result: &pb.Result{Code: 0, Msg: ""}, RowsAffected: rowCount, LastInsertId: lastInsertID, Data: datas}, nil
}
