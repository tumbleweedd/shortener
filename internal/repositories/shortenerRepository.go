package repositories

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/tumbleweedd/shortener/pkg/utils/base62Algorithm"
	"github.com/tumbleweedd/shortener/pkg/utils/sha256Algorithm"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ShortenerRepository struct {
	client *redis.Client
}

func NewShortenerRepository(client *redis.Client) *ShortenerRepository {
	return &ShortenerRepository{client: client}
}

func generateRandomNumber() uint64 {
	rand.Seed(time.Now().UnixNano())

	// Максимальное значение для числа с длиной 8
	maxValue := uint64(99999999)

	return rand.Uint64() % (maxValue + 1)
}

func (r *ShortenerRepository) SaveURL(ctx context.Context, url string) (string, error) {
	code, exists, err := exists(r, url)
	if err != nil {
		log.Println("exists err")
		return "", err
	}

	if exists {
		uintCode, err := strconv.ParseUint(code, 10, 64)
		if err != nil {
			return "", err
		}

		return base62Algorithm.Encode(uintCode), nil
	}

	id := generateRandomNumber()
	key := strconv.FormatUint(id, 10)

	err = r.client.Set(ctx, key, url, 0).Err()
	if err != nil {
		logrus.Println("long url save error")
		return "", err
	}

	hash := sha256Algorithm.CalculateHash(url)
	err = r.client.Set(ctx, hash, key, 0).Err()
	if err != nil {
		log.Println("Hash index save error")
		r.client.Del(ctx, key)
		return "", err
	}

	return base62Algorithm.Encode(id), err
}

func (r *ShortenerRepository) GetURL(ctx context.Context, code string) (string, error) {
	decodedID, err := base62Algorithm.Decode(code)
	if err != nil {
		return "", err
	}

	url, err := r.client.Get(ctx, strconv.FormatUint(decodedID, 10)).Result()
	if err != nil {
		logrus.Println("long URL get error")
		return "", err
	} else if err == redis.Nil {
		return "", errors.New("url error link")
	}

	return url, nil
}

func exists(r *ShortenerRepository, url string) (string, bool, error) {
	hash := sha256Algorithm.CalculateHash(url)
	key, err := r.client.Get(context.Background(), hash).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return key, true, nil
}

/*func (r *ShortenerRepository) SaveURL(ctx context.Context, url string) (string, error) {
	code, exists, err := exists(r, url)
	if err != nil {
		logrus.Println("exists err")
		return "", err
	}

	if exists {
		uintCode, err := strconv.ParseUint(code, 10, 64)
		if err != nil {
			return "", err
		}
		return base62Algorithm.Encode(uintCode), nil
	}

	id := generateRandomNumber()

	res, err := r.client.HSet(ctx, strconv.FormatUint(id, 10), "url", url).Result()
	if err != nil {
		logrus.Println("long url save error")
		return "", err
	}

	fmt.Println(res)

	return base62Algorithm.Encode(id), nil
}
*/

/*func (r *ShortenerRepository) GetURL(ctx context.Context, code string) (string, error) {
	decodedID, err := base62Algorithm.Decode(code)
	if err != nil {
		return "", err
	}

	url, err := r.client.HGet(ctx, strconv.FormatUint(decodedID, 10), "url").Result()
	if err != nil {
		logrus.Println("long URL get error")
		return "", err
	} else if err == redis.Nil {
		return "", errors.New("url error link")
	}

	return url, nil
}

func exists(r *ShortenerRepository, url string) (string, bool, error) {
	keys, err := r.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return "", false, err
	}

	for _, key := range keys {
		values, err := r.client.HMGet(context.Background(), key, "url").Result()
		if err != nil {
			return "", false, err
		}

		if len(values) == 1 && values[0] == url {
			return key, true, nil
		}
	}

	return "", false, nil
}*/
