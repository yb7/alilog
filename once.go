package alilog

import (
  "errors"
  "reflect"
  "sync"
)

//type onceVo struct {
//  Once *sync.Once
//  ExpiresAt time.Time
//  Data interface{}
//}

var onceMap sync.Map
var lockMap sync.Map


func setV(source, dst interface{}) error {
  // ValueOf to enter reflect-land
  dstPtrValue := reflect.ValueOf(dst)
  if dstPtrValue.Kind() != reflect.Ptr {
    return errors.New("destination must be kind of ptr")
  }
  if dstPtrValue.IsNil() {
    return errors.New("destination cannot be nil")
  }
  //dstType := dstPtrType.Elem()
  // the *dst in *dst = zero
  dstValue := reflect.Indirect(dstPtrValue)
  // the = in *dst = 0
  dstValue.Set(reflect.ValueOf(source))
  return nil
}

func loadOnce(key string) *sync.Once {
  lockI, _ := lockMap.LoadOrStore(key, &sync.Mutex{})
  lock := lockI.(*sync.Mutex)

  lock.Lock()

  onceObj, ok := onceMap.Load(key)
  if !ok {
    onceObj = &sync.Once{}
    onceMap.Store(key, onceObj)
  }

  lock.Unlock()
  return onceObj.(*sync.Once)
}
func doOnce(key string, fallback func() error) {
  newOnce := loadOnce(key)

  //var err error
  newOnce.Do(func() {
    //var result interface{}
    if err := fallback(); err != nil {
      onceMap.Delete(key)
    }
    //result, err = fallback()
    //if err == nil {
    //  newOnce.Data = result
    //  onceMap.Store(key, newOnce)
    //}
  })
  //if err != nil {
  //  onceMap.Delete(key)
  //  return err
  //} else {
  //  onceObj, ok := onceMap.Load(key)
  //  if ok {
  //    once := onceObj.(*onceVo)
  //    if once.Data != nil {
  //      setV(once.Data, dst)
  //    } else {
  //      onceMap.Delete(key)
  //    }
  //  }
  //}
  //return nil
}
