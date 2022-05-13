package finder

import (
    "github.com/go-ego/riot"
    "github.com/go-ego/riot/types"
    "os"
    "os/signal"
    "path"
    "sync"
    "syscall"
)

type (
    Engine struct {
    stopCh  chan struct{}
    locker  *Locker
    mu      sync.Mutex
    buckets map[string]*Bucket
    baseDir string
}
    Bucket struct {
        searcher *riot.Engine
    }
)
func NewEngine(baseDir...string) *Engine {
    eng := &Engine{
        stopCh:  make(chan struct{}),
        buckets: map[string]*Bucket{},
        locker:  NewLocker(),
        baseDir: "finder_data",
    }
    if len(baseDir) > 0 {}
    eng.baseDir = baseDir[0]
    return eng
}
func (e *Engine) Bucket(name string) *Bucket {
    e.locker.Lock(name)
    defer e.locker.UnLock(name)
    if b, ok := e.buckets[name]; ok {
        return b
    }
    searcher := &riot.Engine{}

    searcher.Init(types.EngineOpts{
        StoreFolder: path.Join(e.baseDir, name),
        UseStore:    true,
    })
    searcher.FlushIndex()
    b := &Bucket{searcher: searcher}
    e.mu.Lock()
    e.buckets[name] = b
    e.mu.Unlock()
    return b
}
func (e *Engine) WaitStopSignal() {
    handlerArray := []os.Signal{
        syscall.SIGINT,
        syscall.SIGILL,
        syscall.SIGFPE,
        syscall.SIGSEGV,
        syscall.SIGTERM,
        syscall.SIGABRT,
    }
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, handlerArray...)
    select {
    case <-sigs:
        e.onStop()
    case <-e.stopCh:
        e.onStop()
    }
}
func (e *Engine) onStop() {
    e.mu.Lock()
    defer e.mu.Unlock()
    for k, b := range e.buckets {
        b.Close()
        delete(e.buckets, k)
    }
}
func (e *Engine) Stop() {
    e.stopCh <- struct{}{}
}
func (e *Bucket) Index(id string, content string, labels ...string) {
    doc := types.DocData{
        Content: content,
        Attri:   nil,
        Tokens:  nil,
        Labels:  labels,
        Fields:  nil,
    }
    e.searcher.IndexDoc(id, doc)
}
func (e *Bucket) Flush() {
    e.searcher.Flush()
}
func (e *Bucket) Close() {
    e.searcher.Close()
}
func (e *Bucket) Find(content string) (result []string) {
    req := types.SearchReq{
        Text: content,
    }
    resp := e.searcher.Search(req)
    switch val := resp.Docs.(type) {
    case types.ScoredDocs:
        for _, info := range val {
            result = append(result, info.DocId)
        }
    }

    return
}
func (e *Bucket) Remove(id string) {
    e.searcher.RemoveDoc(id)
}
func (e *Bucket) Has(id string) bool {
    return e.searcher.HasDoc(id)
}
