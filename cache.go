package authority

import "time"

var AuthorityModuleIdCache map[int64]*AuthorityModule
var AuthorityModuleNameENCache map[string]*AuthorityModule

var AuthorityAuthorIdCache map[int64]*AuthorityAuthor
var AuthorityAuthorNameENCache map[string]*AuthorityAuthor

const AuthorGarbageCollectorTime = 10 * time.Second
const TimeLongForGCAuthor = 5

func init() {
	AuthorityAuthorIdCache = make(map[int64]*AuthorityAuthor)
	AuthorityAuthorNameENCache = make(map[string]*AuthorityAuthor)
	go AuthorGarbageCollector()
}

func AuthorGarbageCollector() {
	for {
		t := time.Now()
		for _, v := range AuthorityAuthorIdCache {
			if t.Sub(v.Timer).Minutes() > TimeLongForGCAuthor {
				delete(AuthorityAuthorIdCache, v.Id)
				delete(AuthorityAuthorNameENCache, v.NameEN)
			}
		}
		time.Sleep(AuthorGarbageCollectorTime)
	}
}
