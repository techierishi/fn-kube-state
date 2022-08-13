package repository

import (
	"fn-kube-state/pkg/util"
	"log"
)

type DAO interface {
	NewKubeQuery() KubeQuery
}

type dao struct{}

func NewDAO() DAO {
	s := &dao{}
	return s
}

func (d *dao) NewKubeQuery() KubeQuery {
	client, err := util.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	return &kubeQuery{
		client: client,
	}
}
