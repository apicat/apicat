<?php

namespace App\Modules\OperateThrottles;

use Illuminate\Http\Request;

/**
 * 同一IP邮箱注册频率限制
 */
class EmailRegisterThrottle
{
    use Throttle;

    /**
     * 限制周期，单位：秒
     *
     * @var int
     */
    public $decaySeconds = 3600;

    /**
     * 允许尝试的最大次数
     *
     * @var int
     */
    public $maxAttempts = 60;

    /**
     * 通过request生成限制器key
     *
     * @param \Illuminate\Http\Request $request
     * @return string
     */
    public function throttleKey(Request $request)
    {
        return $request->ip() . '|email_register';
    }
}