<?php

namespace App\Modules\OperateThrottles;

use Illuminate\Cache\RateLimiter;
use Illuminate\Http\Request;
use Illuminate\Auth\Events\Lockout;

trait Throttle
{
    /**
     * 检查是否已经多次尝试失败
     *
     * @param  \Illuminate\Http\Request  $request
     * @return bool
     */
    public function hasTooManyAttempts(Request $request)
    {
        return $this->limiter()->tooManyAttempts(
            $this->throttleKey($request), $this->maxAttempts()
        );
    }

    /**
     * 增加尝试计数
     *
     * @param  \Illuminate\Http\Request  $request
     * @return void
     */
    public function incrementAttempts(Request $request)
    {
        $this->limiter()->hit(
            $this->throttleKey($request), $this->decaySeconds()
        );
    }

    /**
     * 清除尝试次数缓存
     *
     * @param  \Illuminate\Http\Request  $request
     * @return void
     */
    public function clearAttempts(Request $request)
    {
        $this->limiter()->clear($this->throttleKey($request));
    }

    /**
     * 锁定对应请求
     *
     * @param  \Illuminate\Http\Request  $request
     * @return void
     */
    public function fireLockoutEvent(Request $request)
    {
        event(new Lockout($request));
    }

    /**
     * 获取限制器实例
     *
     * @return \Illuminate\Cache\RateLimiter
     */
    public function limiter()
    {
        return app(RateLimiter::class);
    }

    /**
     * 获取允许尝试的最大次数
     *
     * @return int
     */
    protected function maxAttempts()
    {
        return property_exists($this, 'maxAttempts') ? $this->maxAttempts : 5;
    }

    /**
     * 获取限制周期，单位：秒
     *
     * @return int
     */
    protected function decaySeconds()
    {
        return property_exists($this, 'decaySeconds') ? $this->decaySeconds : 180;
    }
}
