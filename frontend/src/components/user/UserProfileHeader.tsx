'use client';

import Image from 'next/image';
import { motion } from 'framer-motion';
import { User, Lock, ExternalLink, Verified } from 'lucide-react';
import type { InstagramUser } from '@/types';
import { formatCount } from '@/lib/utils';

interface UserProfileHeaderProps {
  user: InstagramUser;
}

export default function UserProfileHeader({ user }: UserProfileHeaderProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
      className="glass rounded-2xl p-6"
    >
      <div className="flex flex-col sm:flex-row items-center sm:items-start gap-5">
        {/* Profile Picture */}
        <div className="relative h-20 w-20 rounded-full overflow-hidden ring-2 ring-brand-500/30 flex-shrink-0">
          {user.profilePicUrl ? (
            <Image
              src={user.profilePicUrl}
              alt={user.username}
              fill
              className="object-cover"
              sizes="80px"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center bg-surface-700">
              <User className="h-8 w-8 text-surface-500" />
            </div>
          )}
        </div>

        {/* Info */}
        <div className="flex-1 text-center sm:text-left space-y-3">
          <div>
            <div className="flex items-center justify-center sm:justify-start gap-2">
              <h2 className="text-xl font-bold text-surface-100">{user.fullName}</h2>
              {user.isVerified && (
                <Verified className="h-4 w-4 text-brand-400 fill-brand-400/20" />
              )}
            </div>
            <p className="text-sm text-surface-400">@{user.username}</p>
          </div>

          {/* Bio */}
          {user.biography && (
            <p className="text-sm text-surface-300 max-w-lg leading-relaxed">
              {user.biography}
            </p>
          )}

          {/* External URL */}
          {user.externalUrl && (
            <a
              href={user.externalUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-xs text-brand-400 hover:text-brand-300 transition-colors"
            >
              <ExternalLink className="h-3 w-3" />
              {user.externalUrl.replace(/^https?:\/\//, '').split('/')[0]}
            </a>
          )}

          {/* Stats */}
          <div className="flex items-center justify-center sm:justify-start gap-6 pt-1">
            <div className="text-center">
              <p className="text-sm font-bold text-surface-100">
                {formatCount(user.postCount)}
              </p>
              <p className="text-[10px] text-surface-500 uppercase tracking-wider">Posts</p>
            </div>
            <div className="text-center">
              <p className="text-sm font-bold text-surface-100">
                {formatCount(user.followerCount)}
              </p>
              <p className="text-[10px] text-surface-500 uppercase tracking-wider">Followers</p>
            </div>
            <div className="text-center">
              <p className="text-sm font-bold text-surface-100">
                {formatCount(user.followingCount)}
              </p>
              <p className="text-[10px] text-surface-500 uppercase tracking-wider">Following</p>
            </div>
          </div>

          {/* Private badge */}
          {user.isPrivate && (
            <div className="flex items-center justify-center sm:justify-start gap-1.5 text-xs text-yellow-400">
              <Lock className="h-3 w-3" />
              Private account — only public posts are visible
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
}
