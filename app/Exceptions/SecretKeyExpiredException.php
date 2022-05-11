<?php

namespace App\Exceptions;

use Exception;

class SecretKeyExpiredException extends Exception
{
    /**
     * 访问秘钥失效
     *
     * @return array
     */
    public function render()
    {
        return ['status' => -103, 'msg' => '无效的访问秘钥'];
    }
}
