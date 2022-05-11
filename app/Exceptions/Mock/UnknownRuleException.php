<?php

namespace App\Exceptions\Mock;

use Exception;

/**
 * mock语法错误异常
 */
class UnknownRuleException extends Exception
{
    public function render()
    {
        return response()->json([
            'status' => -1,
            'msg' => $this->message . ' 参数mock规则有误，请修改后重试。'
        ]);
    }
}
