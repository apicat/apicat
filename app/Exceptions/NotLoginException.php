<?php

namespace App\Exceptions;

use Exception;
use Illuminate\Http\JsonResponse;

class NotLoginException extends Exception
{
    /**
     * 未登录
     *
     * @return JsonResponse
     */
    public function render()
    {
        return response()->json(['status' => -2, 'msg' => '未登录']);
    }
}
