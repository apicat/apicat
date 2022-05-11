<?php

namespace App\Exceptions\Mock;

use Exception;

/**
 * 参数内容错误异常
 */
class ParamErrorException extends Exception
{
    public function render()
    {
        return response()->json([
            'status' => -1,
            'msg' => 'API返回参数内容有误，请修改后重试。'
        ]);
    }
}
