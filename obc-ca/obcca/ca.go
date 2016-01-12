/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package obcca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/sha3"
)

// CA is the base certificate authority.
//
type CA struct {
	db *sql.DB

	priv *ecdsa.PrivateKey
	cert *x509.Certificate
	raw  []byte
}

// NewCA sets up a new CA.
//
func NewCA(name string) *CA {
	if _, err := os.Stat(RootPath); err != nil {
		Info.Println("Fresh start; creating databases, key pairs, and certificates.")

		if err := os.Mkdir(RootPath, 0755); err != nil {
			Panic.Panicln(err)
		}
	}

	ca := new(CA)

	// open or create certificate database
	db, err := sql.Open("sqlite3", RootPath+"/"+name+".db")
	if err != nil {
		Panic.Panicln(err)
	}

	if err := db.Ping(); err != nil {
		Panic.Panicln(err)
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Certificates (row INTEGER PRIMARY KEY, id VARCHAR(64), timestamp INTEGER, usage INTEGER, cert BLOB, hash BLOB)"); err != nil {
		Panic.Panicln(err)
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users (row INTEGER PRIMARY KEY, id VARCHAR(64), token BLOB, state INTEGER)"); err != nil {
		Panic.Panicln(err)
	}
	ca.db = db

	// read or create signing key pair
	priv, err := ca.readCAPrivateKey(name)
	if err != nil {
		priv = ca.createCAKeyPair(name)
	}
	ca.priv = priv

	// read CA certificate, or create a self-signed CA certificate
	raw, err := ca.readCACertificate(name)
	if err != nil {
		raw = ca.createCACertificate(name, &ca.priv.PublicKey)
	}
	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		Panic.Panicln(err)
	}

	ca.raw = raw
	ca.cert = cert

	return ca
}

// Close closes down the CA.
//
func (ca *CA) Close() {
	ca.db.Close()
}

func (ca *CA) createCAKeyPair(name string) *ecdsa.PrivateKey {
	Trace.Println("Creating CA key pair.")

	curve := elliptic.P384()

	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err == nil {
		raw, _ := x509.MarshalECPrivateKey(priv)
		cooked := pem.EncodeToMemory(
			&pem.Block{
				Type:  "ECDSA PRIVATE KEY",
				Bytes: raw,
			})
		err := ioutil.WriteFile(RootPath+"/"+name+".priv", cooked, 0644)
		if err != nil {
			Panic.Panicln(err)
		}

		raw, _ = x509.MarshalPKIXPublicKey(&priv.PublicKey)
		cooked = pem.EncodeToMemory(
			&pem.Block{
				Type:  "ECDSA PUBLIC KEY",
				Bytes: raw,
			})
		err = ioutil.WriteFile(RootPath+"/"+name+".pub", cooked, 0644)
		if err != nil {
			Panic.Panicln(err)
		}
	}
	if err != nil {
		Panic.Panicln(err)
	}

	return priv
}

func (ca *CA) readCAPrivateKey(name string) (*ecdsa.PrivateKey, error) {
	Trace.Println("Reading CA private key.")

	cooked, err := ioutil.ReadFile(RootPath + "/" + name + ".priv")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(cooked)
	return x509.ParseECPrivateKey(block.Bytes)
}

func (ca *CA) createCACertificate(name string, pub *ecdsa.PublicKey) []byte {
	Trace.Println("Creating CA certificate.")

	raw, err := ca.newCertificate(pub, x509.KeyUsageDigitalSignature, nil)
	if err != nil {
		Panic.Panicln(err)
	}

	cooked := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: raw,
		})
	err = ioutil.WriteFile(RootPath+"/"+name+".cert", cooked, 0644)
	if err != nil {
		Panic.Panicln(err)
	}

	return raw
}

func (ca *CA) readCACertificate(name string) ([]byte, error) {
	Trace.Println("Reading CA certificate.")

	cooked, err := ioutil.ReadFile(RootPath + "/" + name + ".cert")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(cooked)
	return block.Bytes, nil
}

func (ca *CA) createCertificate(id string, pub *ecdsa.PublicKey, usage x509.KeyUsage, timestamp int64, opt ...pkix.Extension) ([]byte, error) {
	Trace.Println("Creating certificate for "+id+".")

	raw, err := ca.newCertificate(pub, usage, opt)
	if err != nil {
		Error.Println(err)
		return nil, err
	}

	hash := sha3.New384()
	hash.Write(raw)
	if _, err = ca.db.Exec("INSERT INTO Certificates (id, timestamp, usage, cert, hash) VALUES (?, ?, ?, ?, ?)", id, timestamp, usage, raw, hash.Sum(nil)); err != nil {
		Error.Println(err)
	}

	return raw, err
}

func (ca *CA) newCertificate(pub *ecdsa.PublicKey, usage x509.KeyUsage, ext []pkix.Extension) ([]byte, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 90)

	parent := ca.cert
	isCA := parent == nil

	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "OBC",
			Organization: []string{"IBM"},
			Country:      []string{"US"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		SubjectKeyId:       []byte{1, 2, 3, 4},
		SignatureAlgorithm: x509.ECDSAWithSHA384,
		KeyUsage:           usage,

		BasicConstraintsValid: true,
		IsCA: isCA,
	}

	if len(ext) > 0 {
		tmpl.Extensions = ext
		tmpl.ExtraExtensions = ext
	}
	if isCA {
		parent = &tmpl
	}

	raw, err := x509.CreateCertificate(
		rand.Reader,
		&tmpl,
		parent,
		pub,
		ca.priv,
	)
	if isCA && err != nil {
		Panic.Panicln(err)
	}

	return raw, err
}

func (ca *CA) readCertificate(id string, usage x509.KeyUsage) ([]byte, error) {
	Trace.Println("Reading certificate for "+id+".")

	var raw []byte
	err := ca.db.QueryRow("SELECT cert FROM Certificates WHERE id=? AND usage=?", id, usage).Scan(&raw)

	return raw, err
}

func (ca *CA) readCertificate1(id string, ts int64) ([]byte, error) {
	Trace.Println("Reading certificate for "+id+".")

	var raw []byte
	err := ca.db.QueryRow("SELECT cert FROM Certificates WHERE id=? AND timestamp=?", id, ts).Scan(&raw)

	return raw, err
}

func (ca *CA) readCertificates(id string, opt ...int64) (*sql.Rows, error) {
	Trace.Println("Reading certificatess for "+id+".")

	if len(opt) > 0 && opt[0] != 0 {
		return ca.db.Query("SELECT cert FROM Certificates ORDER BY usage WHERE id=? AND timestamp=?", id, opt[0])
	}

	return ca.db.Query("SELECT cert FROM Certificates WHERE id=?", id)
}

func (ca *CA) readCertificateByHash(hash []byte) ([]byte, error) {
	Trace.Println("Reading certificate for hash "+string(hash)+".")

	var raw []byte
	row := ca.db.QueryRow("SELECT cert FROM Certificates WHERE hash=?", hash)
	err := row.Scan(&raw)

	return raw, err
}

func (ca *CA) registerUser(id string, opt ...string) (string, error) {
	Trace.Println("Registering user "+id+".")

	var row int
	err := ca.db.QueryRow("SELECT row FROM Users WHERE id=?", id).Scan(&row)
	if err == nil {
		return "", errors.New("user is already registered")
	}

	var tok string
	if len(opt) > 0 && len(opt[0]) > 0 {
		tok = opt[0]
	} else {
		tok = randomString(12)
	}

	_, err = ca.db.Exec("INSERT INTO Users (id, token, state) VALUES (?, ?, ?)", id, tok, 0)
	if err != nil {
		Error.Println(err)

	}

	return tok, err
}

func (ca *CA) deleteUser(id string) (error) {
	Trace.Println("Deleting user "+id+".")

	var row int
	err := ca.db.QueryRow("SELECT row FROM Users WHERE id=?", id).Scan(&row)
	if err == nil {
		_, err = ca.db.Exec("DELETE FROM Certificates Where id=?", id)
		if err != nil {
			Error.Println(err)
		}

		_, err = ca.db.Exec("DELETE FROM Users WHERE row=?", row)
		if err != nil {
			Error.Println(err)
		}
	}
	
	return err
}

func (ca *CA) readToken(id string) *sql.Row {
	Trace.Println("Reading token for "+id+".")

	return ca.db.QueryRow("SELECT token, state FROM Users WHERE id=?", id)
}
