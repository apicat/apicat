<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class DocShare extends Model
{
    use HasFactory;

    protected $fillable = [
        'project_id',
        'user_id',
        'doc_id',
        'secret_key',
    ];
}
