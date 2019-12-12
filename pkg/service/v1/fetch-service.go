package v1

import (
	v1 "axxon/pkg/api/v1"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type fetchServiceServer struct {
	db *sql.DB
}

// NewFetchServiceServer creates fetch server
func NewFetchServiceServer(db *sql.DB) v1.FetchServiceServer {
	return &fetchServiceServer{db: db}
}

func (s *fetchServiceServer) checkAPI(api string) error {
	// "" == use current api version
	if len(api) > 0 {
		if api != apiVersion {
			return status.Error(codes.Unimplemented, "unsupported API version")
		}
	}
	return nil
}

func (s *fetchServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database: "+err.Error())
	}
	return c, nil
}

func (s *fetchServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(apiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if req.Request == nil {
		return nil, status.Error(codes.Unknown, "request parameter is required")
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	request, err := http.NewRequest(
		req.Request.Method,
		req.Request.Address,
		nil,
	)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to create request: "+err.Error())
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to send request: "+err.Error())
	}
	defer resp.Body.Close()

	res, err := c.ExecContext(ctx, "INSERT INTO Request(`Method`, `Address`, `Headers`, `Body`) VALUES(?, ?, ?, ?)", req.Request.Method, req.Request.Address, req.Request.Headers, req.Request.Body)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed insert values into request: "+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve last insert id: "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (s *fetchServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(apiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Method`, `Address`, `Headers`, `Body` FROM `Request`")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from request: "+err.Error())
	}
	defer rows.Close()

	list := []*v1.Request{}
	for rows.Next() {
		request := new(v1.Request)
		if err := rows.Scan(&request.Id, &request.Method, &request.Address, &request.Headers, &request.Body); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve values from request: "+err.Error())
		}
		list = append(list, request)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from request: "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:      apiVersion,
		Requests: list,
	}, nil
}

func (s *fetchServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(apiVersion); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM Request WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete: "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected"+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Unknown row with ID=%d", req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}
