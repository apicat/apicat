<?php

namespace App\Repositories\ApiDoc;

use App\Models\ApiDoc;
use App\Models\ApiDocHistories;
use App\Models\ApiDocHistory;
use App\Repositories\Project\TreeCacheRepository;
use Illuminate\Support\Facades\Auth;

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

    /**
     * 历史记录列表
     *
     * @param int $docId 文档id
     * @return \Illuminate\Database\Eloquent\Collection
     */
    public static function list($docId)
    {
        return ApiDocHistories::where('doc_id', $docId)->latest()->get();
    }

    /**
     * 历史记录
     *
     * @param int $id 记录id
     * @return ApiDocHistory
     */
    public static function get(int $id)
    {
        return ApiDocHistories::find($id);
    }

    /**
     * 恢复历史记录
     *
     * @param ApiDocHistories|int $history 历史记录
     * @param ApiDoc $now 当前文档
     * @return void
     */
    public static function restore($history, $now)
    {
        if (is_numeric($history)) {
            if (!$history = ApiDocHistories::find($history)) {
                return;
            }
        }
        if (is_numeric($now)) {
            if (!$now = ApiDoc::find($now)) {
                return;
            }
        }

        ApiDocHistories::create([
            'doc_id' => $now->id,
            'title' => $now->title,
            'content' => $now->content,
            'last_user_id' => $now->updated_user_id,
            'last_updated_at' => $now->updated_at
        ]);

        $now->title = $history->title;
        $now->content = $history->content;
        $now->updated_user_id = Auth::id();
        $now->save();

        TreeCacheRepository::remove($now->project_id);
    }
}
