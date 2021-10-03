package vshttp

import (
	"context"
	"errors"
	"log"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"golang.org/x/crypto/acme/autocert"
)

var ErrCacheMiss = errors.New("acme/autocert: certificate cache miss")

type CertsCache struct{}

// Get reads a certificate data from DB by name
func (d CertsCache) Get(ctx context.Context, name string) ([]byte, error) {
	log.Println("Get crt from cache", name)
	var (
		data []byte
		err  error
		done = make(chan struct{})
	)
	go func() {
		defer close(done)
		var crt models.Crt
		err = db.GetConnection().Where(&models.Crt{Name: name}).First(&crt)
		if err != nil {
			log.Println("err != nil")
			if db.ErrRecordNotFound(err) {
				log.Println("cert not found in DB")
				err = autocert.ErrCacheMiss
			} else {
				log.Println("failed to get crt data cache", err) // nil, ErrCacheMiss
				return
			}
		}
		data = crt.Data
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}
	log.Println("return ", err, data)
	return data, err
}

// Put inserts the certificate data to DB
func (d CertsCache) Put(ctx context.Context, name string, data []byte) error {
	log.Println("Put crt", name)
	done := make(chan struct{})
	var err error
	go func() {
		defer close(done)
		err = db.GetConnection().Save(&models.Crt{
			Name: name,
			Data: data,
		})
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	return err
}

// Delete removes the specified file name.
func (d CertsCache) Delete(ctx context.Context, name string) error {
	var (
		err  error
		done = make(chan struct{})
	)
	go func() {
		defer close(done)
		err = db.GetConnection().Delete(&models.Crt{Name: name})
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	return err
}
