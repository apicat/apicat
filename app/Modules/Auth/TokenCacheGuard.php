<?php

namespace App\Modules\Auth;

use Illuminate\Auth\GuardHelpers;
use Illuminate\Contracts\Auth\Guard;
use Illuminate\Contracts\Auth\UserProvider;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis as LaravelRedis;
use Illuminate\Support\Str;
use App\Models\User;

/**
 * API调用的用户守卫
 */
class TokenCacheGuard implements Guard
{
    use GuardHelpers;

    /**
     * The request instance.
     *
     * @var \Illuminate\Http\Request
     */
    protected $request;

    /**
     * The name of the query string item from the request containing the API token.
     *
     * @var string
     */
    protected $inputKey;

    /**
     * Indicates if the API token is hashed in storage.
     *
     * @var bool
     */
    protected $hash = false;

    /**
     * Redis连接实例
     *
     * @var \Illuminate\Redis\Connections\Connection
     */
    protected $redis;

    /**
     * Redis连接名称，在config/database.php配置文件中设置
     *
     * @var string
     */
    private $redisConnection = 'api_token';

    /**
     * token缓存天数
     *
     * @var integer
     */
    private $expireDay = 7;

    /**
     * Create a new authentication guard.
     *
     * @param \Illuminate\Contracts\Auth\UserProvider $provider
     * @param \Illuminate\Http\Request $request
     * @param string $inputKey
     * @param bool $hash
     * @return void
     */
    public function __construct(UserProvider $provider, Request $request, $inputKey = 'api_token', $hash = true)
    {
        $this->provider = $provider;
        $this->request = $request;
        $this->inputKey = $inputKey;
        $this->hash = $hash;
    }

    /**
     * Get the currently authenticated user.
     *
     * @return \Illuminate\Contracts\Auth\Authenticatable|null
     */
    public function user()
    {
        // If we've already retrieved the user for the current request we can just
        // return it back immediately. We do not want to fetch the user data on
        // every call to this method because that would be tremendously slow.
        if (!is_null($this->user)) {
            return $this->user;
        }

        $user = null;

        $token = $this->getTokenForRequest();

        if (!empty($token)) {
            $storageKey = $this->hash ? hash('sha256', $token) : $token;
            $cacheData = $this->getDataFromCache($storageKey);

            if ($cacheData and (isset($cacheData['uid']) or isset($cacheData['id']))) {
                $uid = isset($cacheData['uid']) ? $cacheData['uid'] : $cacheData['id'];
                $user = $this->provider->retrieveById($uid);

                // 延长缓存时间
                $this->putDataToCache($storageKey, $this->cacheUserData($user));
            }
        }

        return $this->user = $user;
    }

    /**
     * Get the token for the current request.
     *
     * @return string
     */
    public function getTokenForRequest()
    {
        $token = $this->request->query($this->inputKey);

        if (empty($token)) {
            $token = $this->request->input($this->inputKey);
        }

        if (empty($token)) {
            $token = $this->request->bearerToken();
        }

        return $token;
    }

    /**
     * Validate a user's credentials.
     *
     * @param array $credentials
     * @return bool
     */
    public function validate(array $credentials = [])
    {
        if (!isset($credentials[$this->inputKey]) or empty($credentials[$this->inputKey])) {
            return false;
        }

        $storageKey = $this->hash ? hash('sha256', $credentials[$this->inputKey]) : $credentials[$this->inputKey];
        $cacheData = $this->getDataFromCache($storageKey);

        if ($cacheData and (isset($cacheData['uid']) or isset($cacheData['id']))) {
            $uid = isset($cacheData['uid']) ? $cacheData['uid'] : $cacheData['id'];
            if ($this->provider->retrieveById($uid)) {
                return true;
            }
        }

        return false;
    }

    /**
     * 将给定的用户登录到api令牌守卫中
     *
     * @param User $user 用户信息实例
     * @return string
     */
    public function login(User $user)
    {
        $token = Str::random(60);
        $storageKey = $this->hash ? hash('sha256', $token) : $token;
        if (!$this->getDataFromCache($storageKey)) {
            if ($this->putDataToCache($storageKey, $this->cacheUserData($user))) {
                return 'Bearer ' . $token;
            }
        }
    }

    /**
     * 将给定的api token从api令牌守卫中退出
     *
     * @return void
     */
    public function logout()
    {
        $token = $this->getTokenForRequest();
        if (!empty($token)) {
            $storageKey = $this->hash ? hash('sha256', $token) : $token;
            $this->delDataFromCache($storageKey);
        }
    }

    /**
     * 从Redis里获取token数据
     *
     * @param string $cacheKey 缓存键名
     * @return array|null
     */
    protected function getDataFromCache($cacheKey)
    {
        if (!$this->redis) {
            $this->redis = LaravelRedis::connection($this->redisConnection);
        }

        if ($value = $this->redis->get($cacheKey)) {
            if ($data = json_decode($value, true)) {
                return $data;
            }
        }
    }

    /**
     * 向Redis里写入最新的token数据
     *
     * @param string $cacheKey
     * @param string $cacheData
     * @return boolean
     */
    protected function putDataToCache(string $cacheKey, string $cacheData)
    {
        if (!$this->redis) {
            $this->redis = LaravelRedis::connection($this->redisConnection);
        }

        return $this->redis->set($cacheKey, $cacheData, 'EX', $this->expireDay * 86400);
    }

    /**
     * 从Redis里清除token数据
     *
     * @param string $cacheKey 缓存键名
     * @return boolean
     */
    protected function delDataFromCache($cacheKey)
    {
        if (!$this->redis) {
            $this->redis = LaravelRedis::connection($this->redisConnection);
        }

        return $this->redis->del($cacheKey);
    }

    /**
     * 缓存用户数据
     *
     * @param User $user 用户信息实例
     * @return string
     */
    protected function cacheUserData(User $user)
    {
        $data = [
            'uid' => $user->id,
        ];

        return json_encode($data);
    }
}