package query

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/eco/longy/x/longy/internal/querier"
	"github.com/patrickmn/go-cache"
	"net/http"
	"sync"
	"time"
)

var boardCache *cache.Cache
var mutex *sync.Mutex

const (
	//BoardCacheKey is the key for the cached leaderboard
	BoardCacheKey = "leaderboard"
)

func init() {
	boardCache = cache.New(1*time.Minute, 2*time.Minute)
	mutex = &sync.Mutex{}
}

//LeaderBoardHandler queries the chain for for the current leader board
//nolint:gocritic
func LeaderBoardHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, found := boardCache.Get(BoardCacheKey)
		if found {
			rest.PostProcessResponse(w, cliCtx, res)
			return
		}

		//update cache
		mutex.Lock()
		//pass through for any routines that queued on the lock while the cache was updated
		res, found = boardCache.Get(BoardCacheKey)
		if found {
			mutex.Unlock()
			rest.PostProcessResponse(w, cliCtx, res)
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			storeName, querier.LeaderKey), nil)
		if err != nil {
			mutex.Unlock()
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		boardCache.Set(BoardCacheKey, res, cache.DefaultExpiration)

		mutex.Unlock()
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
