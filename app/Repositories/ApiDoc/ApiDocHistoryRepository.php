<?php

namespace App\Repositories\ApiDoc;

use App\Models\ApiDocHistories;
use App\Models\ApiDocHistory;

class ApiDocHistoryRepository
{
    /**
     * 添加历史记录
     *
     * @param int $docId 文档id
     * @param string $title 文档标题
     * @param string $content 文档内容
     * @param int $userID 用户id
     * @param \Illuminate\Support\Carbon $updated_time 编辑时间
     * @return ApiDocHistory
     */
    public static function add($docId, $title, $content, $userID, $updated_time)
    {
        return ApiDocHistories::create([
            'doc_id' => $docId,
            'title' => $title,
            'content' => $content,
            'last_user_id' => $userID,
            'last_updated_at' => $updated_time
        ]);
    }
}
