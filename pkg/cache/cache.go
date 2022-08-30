package cache

import (
	"errors"
	"mytasks"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[int64]Item
	database          *sqlx.DB
}

func NewCache(db *sqlx.DB, defaultExpiration, cleanupInterval time.Duration) *Cache {
	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[int64]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC() // данный метод рассматривается ниже
	}

	cache.database = db

	return &cache
}

func (c *Cache) GetItems() map[int64]Item {
	c.Lock()
	defer c.Unlock()
	return c.items
}

func (c *Cache) GetDefaultExpiration() time.Duration {
	return c.defaultExpiration
}

func (c *Cache) Set(key int64, value interface{}, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}
func (c *Cache) GetFromDB() ([]mytasks.TaskList, error) {
	var lists []mytasks.TaskList
	query := "SELECT * FROM tasks"
	err := c.database.Select(&lists, query)
	return lists, err
}

func (c *Cache) FullToRam() error {
	lists, err := c.GetFromDB()
	if err != nil {
		return err
	}
	for _, list := range lists {
		c.Set(list.ID, list, viper.GetDuration("cache.expiration"))
	}
	return nil
}

func (c *Cache) Refresh(refreshTime time.Duration) {
	err := c.FullToRam()
	if err != nil {
		logrus.Error("Problems with cache: %s", err.Error())
	}
	//	fmt.Println("Updating cache")
	//	fmt.Println(fmt.Sprintf("adress cache: %v \n adress value: %v", &c, c))
	ticker := time.NewTicker(refreshTime * time.Millisecond)
	for {
		select {
		case _ = <-ticker.C:
			err := c.FullToRam()
			if err != nil {
				logrus.Error("Problems with cache: %s", err.Error())
			}
			//			fmt.Println("Updating cache")
			//			fmt.Println(len(c.items))
		}
	}
}

func (c *Cache) Get(key int64) (interface{}, bool) {

	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check - на установку времени истечения
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

func (c *Cache) Delete(key int64) error {

	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}
	delete(c.items, key)
	return nil
}

///
func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []int64) {

	c.RLock()
	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []int64) {

	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
