package main

import (
	"errors"

	"gorm.io/gorm"
)

type TLDRDBCached struct {
	db       *gorm.DB
	provider TLDRProvider
}

func (t *TLDRDBCached) Retrieve(key string) string {
	var ent TLDREntity
	err := t.db.First(&ent, "Key = ?", key).Error
	switch {
	case err == nil:
		return ent.Val
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return t.provider.Retrieve(key)

	}
	val := t.provider.Retrieve(key)
	_ = t.db.Create(&TLDREntity{Key: key, Val: val}).Error
	return val
}

func (t *TLDRDBCached) List() []string {
	all := t.provider.List()
	if len(all) == 0 {
		return nil
	}
	var cachedRecords []TLDREntity
	t.db.Select("Key").Find(&cachedRecords)
	cachedSet := make(map[string]struct{}, len(cachedRecords))
	for _, rec := range cachedRecords {
		cachedSet[rec.Key] = struct{}{}
	}
	var inCache, notInCache []string
	for _, key := range all {
		if _, ok := cachedSet[key]; ok {
			inCache = append(inCache, key)
		} else {
			notInCache = append(notInCache, key)
		}
	}
	return append(inCache, notInCache...)
}

func NewTLDRDBCached(nonCachedProvider TLDRProvider) TLDRProvider {
	return &TLDRDBCached{
		db:       GetConnection(),
		provider: nonCachedProvider,
	}
}
