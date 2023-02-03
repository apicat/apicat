<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ApiDocHistories extends Model
{
    use HasFactory;

    protected $fillable = [
        'doc_id', 'title', 'content', 'last_user_id', 'last_updated_at'
    ];

    protected $dates = [
        'last_updated_at'
    ];
}
