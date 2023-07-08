package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("exclude first added if overloaded", func(t *testing.T) {
		c := NewCache(3)
		// Заполняем кэш с переполнением
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)

		// Устанавливаем значение первому добавленному ключу
		wasInCache := c.Set("aaa", 111)

		require.False(t, wasInCache)
	})
	t.Run("exclude long time unused if overloaded", func(t *testing.T) {
		c := NewCache(3)

		// Заполняем кэш
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		// Работаем с ранее добавленными элементами
		c.Set("ccc", 333)
		c.Get("aaa")
		c.Get("bbb")

		// Добавляем новое значение (превышаем capacity)
		c.Set("ddd", 400)

		// Устанавливаем значение последнему добавленному ключу
		wasInCache := c.Set("ccc", 123)

		require.False(t, wasInCache)
	})
	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(3)
		var ok bool

		// Заполняем кэш
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		// Проверяем, что значения лежат в кэше
		_, ok = c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("bbb")
		require.True(t, ok)
		_, ok = c.Get("ccc")
		require.True(t, ok)

		c.Clear() // Очищаем кэш

		// Проверяем, что значений в кэше нет
		_, ok = c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
		_, ok = c.Get("ccc")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
