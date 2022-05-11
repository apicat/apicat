<?php

namespace App\Exceptions\Mock;

use Exception;

/**
 * 无返回参数异常
 */
class NoParamException extends Exception
{
    public function render()
    {
        return response()->json([
            'status' => -1,
            'msg' => '未添加API返回参数，请添加后重试。'
        ]);
    }
}
